from flask import Flask, render_template, request
import requests

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

# Flask路由
@app.route('/')
def index():
    return render_template('index.html')

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