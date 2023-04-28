import datetime

import requests
from flask import Flask, jsonify, render_template, request

app = Flask(__name__, template_folder='./')

# 限制单IP每日10次
MAX_REQUESTS_PER_DAY = 10

# 用于生成API签名的密钥对
API_KEY = 'xxxxxxxxx'
API_SECRET = 'xxxxxxxxxxx'

# 缓存刷新API的URL
API_URL = 'https://xxxxxxxxxxxx/v1/jobs'

# 当前IP已经发出的请求数量
request_count = {}


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
        if not url.startswith('https://github.com'):
            return 'Invalid URL 2'
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
        if not url.startswith('https://www.npmjs.com/package/'):
            return 'Invalid URL 2'
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

# 刷新CDN缓存
def refresh_cdn_cache(urls, is_dir):
    headers = {
        'api-key': API_KEY,
        'api-secret': API_SECRET
    }

    json_data = []
    for url in urls:
        data = {
            "type": "clean_url" if not is_dir else "clean_dir",
            "data": {
                "url": url
            }
        }
        json_data.append(data)

    response = requests.post(API_URL, json=json_data, headers=headers)

    # 获取响应结果并返回
    result = response.json()
    return result.get('code'), result.get('msg')

@app.route('/refresh', methods=['POST'])
def refresh():
    global request_count

    # 获取提交的数据
    urls = request.form.getlist('url')
    is_dir = 'is_dir' in request.form

    # 检查请求次数是否超过限制
    remote_addr = request.remote_addr
    if remote_addr in request_count:
        if request_count[remote_addr] >= MAX_REQUESTS_PER_DAY:
            return '请求次数超过限制', 429
        else:
            request_count[remote_addr] += 1
    else:
        request_count[remote_addr] = 1

    # 刷新CDN缓存并返回结果
    code, msg = refresh_cdn_cache(urls, is_dir)
    return f'刷新结果：{msg}（code={code}）'

if __name__ == '__main__':
    app.run(debug=True)
