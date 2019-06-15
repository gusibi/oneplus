#! -*- coding: utf-8 -*-

# 1. 读取markdown 中code
# 2. 把code 生成html
# 3. html 转成图片
# 4. 上传图片到七牛
# 5. 图片地址替换到原markdown


class OnePlush:

    def __init__(self, mfile_path, theme="atom-one-dark", language="c"):
        self.mfile_path = mfile_path
        self.theme = theme
        self.language = language

    def read_markdown(self):
        pass

    def find_codes(self):
        pass

    def generate_html(self, codes):
        for code in codes:
            pass

    def html_to_img(self):
        pass

    def upload_img(self):
        pass

    def update_markdown(self):
        pass


def main():
    pass
