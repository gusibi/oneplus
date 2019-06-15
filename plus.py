#! -*- coding: utf-8 -*-

import re
import os
import sys
import codecs
import logging
from os import environ
from optparse import OptionParser

# 1. 读取markdown 中code
# 2. 把code 生成html
# 3. html 转成图片
# 4. 上传图片到七牛
# 5. 图片地址替换到原markdown

try:
    import hashlib

    md5_constructor = hashlib.md5
    md5_hmac = md5_constructor
    sha_constructor = hashlib.sha1
    sha_hmac = sha_constructor
except ImportError:
    import md5

    md5_constructor = md5.new
    md5_hmac = md5
    import sha

    sha_constructor = sha.new
    sha_hmac = sha

import requests
from qiniu import Auth, put_file, Zone, set_default

from webshot import webshot

Code2HtmlCreateUrl = "http://service-g5235zgh-1254035985.ap-beijing.apigateway.myqcloud.com/test/codes"
Code2HtmlUrlPrefix = "http://service-g5235zgh-1254035985.ap-beijing.apigateway.myqcloud.com/test/codes/"

logger = logging.getLogger("oneplus: ")
# formatter = logging.Formatter("%(asctime)s - %(name)s - %(levelname)s - %(message)s")
# logger.setFormatter(formatter)

QINIU_BUCKET_NAME = "blog"
QINIU_ACCESS_TOKEN = environ.get("QINIU_ACCESS_TOKEN")
QINIU_SECRET_TOKEN = environ.get("QINIU_SECRET_TOKEN")
QINIU_DOMAIN = "media.gusibi.mobi"

def MD5(code):
    return md5_constructor(code.encode("utf-8")).hexdigest().lower()

def upload27niu(filepath):
    if not filepath:
        logger.error("filepath is invalid")
        return
    key = filepath.split("/")[-1]
    logger.info("upload image: %s to 7niu" % key)
    access_key = str(QINIU_ACCESS_TOKEN)
    secret_key = str(QINIU_SECRET_TOKEN)
    q = Auth(access_key, secret_key)
    token = q.upload_token(QINIU_BUCKET_NAME, key, 360)
    ret, info = put_file(token, key, filepath)
    error_types = ["text/plain"]
    if info.status_code == 200:
        if ret.get("mimeType") in error_types:
            # 如果上传错误
            logger.error("image upload error")
            return None, True
        url = "http://%s/%s" % (QINIU_DOMAIN, key) 
        logger.info("url: %s" % url)
        return url, False
    return None, True

def code2html(code, language="plaintext"):
    data = {
        "code": code,
        "language": language
    }
    resp = requests.post(Code2HtmlCreateUrl, json=data)
    if resp.status_code == 200:
        content = resp.json()
    else:
        content = resp.content
    return content["ID"]

def code2img(code, language):
    cid = MD5(code)
    outfile = "/tmp/%s.png" % cid
    if os.path.exists(outfile):
        return outfile
    cid = code2html(code, language)
    html_url = "%s%s" % (Code2HtmlUrlPrefix, cid)
    lines = len(code.split("\n"))
    webshot(html_url, 20*lines, outfile)
    return outfile


class OnePlush:

    _fenced_code_block_re = re.compile(r'''
        (?:\n+|\A\n?)
        ^```\s*?([\w+-]+)?\s*?\n    # opening fence, $1 = optional lang
        (.*?)                       # $2 = code block content
        ^```[ \t]*\n                # closing fence
        ''', re.M | re.X | re.S)

    def __init__(self, mfile_path, new_filepath, theme="atom-one-dark"):
        self.mfile_path = mfile_path
        self.new_filepath = new_filepath
        self.theme = theme
        self.text = self.read_markdown(self.mfile_path)
        self.save_new_markdown(self.text)

    def read_markdown(self, path, encoding="utf-8"):
        fp = codecs.open(path, 'r', encoding)
        text = fp.read()
        fp.close()
        return "%s\n" % text

    def _fenced_code_block_sub(self, match):
        language = match.group(1)
        codeblock = match.group(2)
        image_path = code2img(codeblock, language)
        image_url, _ = upload27niu(image_path)
        return "\n\n![](%s)\n\n" % image_url
        # return self._code_block_sub(match, is_fenced_code_block=True)

    def _do_fenced_code_blocks(self, text):
        """Process ```-fenced unindented code blocks ('fenced-code-blocks' extra)."""
        return self._fenced_code_block_re.sub(self._fenced_code_block_sub, text)

    def save_new_markdown(self, text):
        # codes = self._fenced_code_block_re.findall(text)
        # for language, code in codes:
        #     print(language, code)
        new_text = self._do_fenced_code_blocks(text)
        with open(self.new_filepath, 'w') as f:
            f.write(new_text.encode("utf-8"))
            print("save new markdown to %s" % self.new_filepath)


def main():
    parser = OptionParser(usage="", add_help_option=True)
    parser.add_option("-m", "--markdown", nargs=1, type="string", dest="markdown",
                      help="markdown", metavar="FILE")
    parser.add_option("-n", "--newfile", nargs=1, dest="newfile", type="string",
                      help="outfile", metavar="FILE")
    (options, args) = parser.parse_args()
    markdown = options.markdown
    if not os.path.exists(markdown):
        print("No such file or directory: %s" % markdown)
        return markdown
    print(options, markdown, options.newfile)
    OnePlush(options.markdown, options.newfile)


if __name__ == "__main__":
    main()
