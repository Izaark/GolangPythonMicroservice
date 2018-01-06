from flask import Flask, render_template, request, Response
import requests

app = Flask(__name__)

@app.route('/', methods = ['GET'])
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

@app.route('/pokemon')
def index():
	return render_template('index.html')

if __name__ == '__main__':
	app.run(debug = True)