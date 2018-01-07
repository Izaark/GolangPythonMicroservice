from flask import Flask, render_template, request, Response
import requests, json

app = Flask(__name__)

@app.route('/pokemon', methods = ['GET'])
def get_pokemon():

	urlApi = "http://localhost:1500/api/pokemon"
	r = requests.get(urlApi)
	json_body = r.json()
	
	message = json_body["message"]
	pokemon_name = json_body['response']

	for poke in pokemon_name:
		pokname = poke["name"]
	#return Response("{'a':'b'}", status=200, mimetype='application/json')
	return render_template('pokemon.html', poke_name=pokemon_name), 200


@app.route('/register/pokemon', methods = ['POST'])
def post_pokemon():
	pname = request.form['pokemon']
	purl = request.form['url']

	urlApi = "http://localhost:1500/api/register/pokemon"
	payload = {'name':pname, 'url':purl}
	#headers = {'Content-Type':'application/json'}
	r = requests.post(urlApi, data=json.dumps(payload))

	if r.status_code == 200:
		headers_response = r.headers
		print(headers_response)
		body = r.json()
		print(body)
		return render_template('new_pokemon.html', npokemon=pname, nurl=purl), 200
	elif r.status_code == 409:
		body = r.json()
		message = body['message']
		status = body['status']
		return render_template('http_code.html', nerror= message, nstatus=status), 409

#todo: function: Delete, Update, GetAPokemon
@app.route('/')
def index():
	return render_template('index.html')

if __name__ == '__main__':
	app.run(debug = True)