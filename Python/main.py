from flask import Flask, render_template, request, Response
import requests, json

app = Flask(__name__)

@app.route('/pokemon', methods = ['GET'])
def get_pokemon():

	#todo set permission ports for sending to Api
	urlApi = "http://localhost:1500/api/pokemon"
	r = requests.get(urlApi)
	json_body = r.json()
	if r.status_code == 200:
		message = json_body["message"]
		pokemon_name = json_body['response']

		for poke in pokemon_name:
			pokname = poke["name"]
		#return Response("{'a':'b'}", status=200, mimetype='application/json')
		return render_template('pokemon.html', poke_name=pokemon_name), 200
	else:
		print('error')

#todo egenrate models !
@app.route('/register/pokemon', methods = ['POST'])
def post_pokemon():
	pname = request.form['pokemon']
	purl = request.form['url']
	#todo set permission ports for sending to Api
	if request.method == 'POST':		
		urlApi = "http://localhost:1500/api/register/pokemon"
		payload = {'name':pname, 'url':purl}
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

@app.route('/update/pokemon/<id>', methods = ['PUT', 'DELETE'])
def update_pokemon(id):
	print(id)
	if request.method == 'PUT':
		pname = request.form['pokemon']
		purl = request.form['url']

		urlApi = 'http://localhost:1500/api/update/pokemon/'+id
		payload = {'name':pname, 'url':purl}
		r = requests.put(urlApi, data=json.dumps(payload))
		if r.status_code == 200 :
			body = r.json()
			print(body)
			return render_template('new_pokemon.html', npokemon=pname, nurl=purl), 200
		else:
			body = r.json()
			message = body['message']
			status = body['status']
			print(message)
			return render_template('http_code.html', nerror= message, nstatus=status), 500

	elif request.method == 'DELETE':
		urlApi = 'http://localhost:1500/api/delete/pokemon/'+id
		r = requests.delete(urlApi)
		return render_template('delete.html', pid= id), 200

#todo: function: GetAPokemon and more validation !! in funtion
@app.route('/')
def index():
	return render_template('index.html')

if __name__ == '__main__':
	app.run(debug = True,host='0.0.0.0')