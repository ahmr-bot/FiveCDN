import datetime
import requests
from flask import Flask, render_template, request, jsonify

app = Flask(__name__, template_folder='./')


def write_to_file(file_path, content):
    with open(file_path, 'a+') as f:
        f.write(content)


def get_github_api(owner, repo):
    api_url = f'https://api.github.com/repos/{owner}/{repo}'
    headers = {'Accept': 'application/vnd.github.v3+json'}
    response = requests.get(api_url, headers=headers)
    if response.status_code != 200:
        return None
    return response.json()


@app.route('/')
def index():
    return render_template('index.html')


@app.route('/check')
def check():
    return render_template('check.html')


@app.route('/submit', methods=['POST'])
def submit():
    platform = request.args.get('platform')
    url = request.args.get('url')
    # url = 'https://github.com/ahmr-bot/Mirrors'
    if platform == 'github':
        split_url = url.split('/')
        owner = split_url[-2]
        repo = split_url[-1]
        with open('pass.txt') as f:
            pass_urls = f.read().splitlines()
        with open('audit.txt') as f:
            audit_urls = f.read().splitlines()
        if any(f'/gh/{owner}/{repo}/' in url for url in pass_urls):
            return 'Already pass'
        if any(f'/gh/{owner}/{repo}/' in url for url in audit_urls):
            return 'Already In Auditing'

        json_response = get_github_api(owner, repo)
        if not json_response:
            return 'Invalid URL 1'

        stargazers_count = json_response['stargazers_count']
        if stargazers_count > 1000:
            write_to_file('pass.txt', f'/gh/{owner}/{repo}/*\n')
        else:
            write_to_file('audit.txt', f'/gh/{owner}/{repo}/*\n')
    elif platform == 'npm':

        package_name = url.replace('https://www.npmjs.com/package/', '')
        with open('pass.txt') as f:
            pass_urls = f.read().splitlines()
        with open('audit.txt') as f:
            audit_urls = f.read().splitlines()
        if any(f'/npm/{package_name}' in url for url in pass_urls):
            return 'Already pass'
        if any(f'/npm/{package_name}' in url for url in audit_urls):
            return 'Already In Auditing'

        now = datetime.date.today()
        two_days_ago = now - datetime.timedelta(days=2)
        today = two_days_ago.strftime('%Y-%m-%d')
        api_url = f'https://npm-stat.com/api/download-counts?package={package_name}&from={today}&until={today}'
        response = requests.get(api_url)

        json_response = response.json()
        if not json_response:
            return 'Invalid URL 2'
        if package_name in json_response and today in json_response[package_name]:
            downloads = json_response[package_name][today]
        else:
            return json_response

        if downloads > 1000:
            write_to_file('pass.txt', f'/npm/{package_name}*\n')
        else:
            write_to_file('audit.txt', f'/npm/{package_name}*\n')

    return 'Success'


@app.route('/api/v1/pass')
def get_pass():
    with open('pass.txt') as f:
        pass_urls = f.read().splitlines()

    return jsonify(urls=pass_urls)


@app.route('/api/v1/audit')
def get_audit():
    with open('audit.txt') as f:
        audit_urls = f.read().splitlines()

    return jsonify(urls=audit_urls)


if __name__ == '__main__':
    app.run(debug=True)
