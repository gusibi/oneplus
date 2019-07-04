# -*- coding: utf-8 -*-

from selenium import webdriver
import time
import os.path
import multiprocessing as mp

def webshot(url, height, outfile):
    driver = webdriver.PhantomJS()
    driver.set_window_size(660, height)
    # driver.maximize_window()
    # 返回网页的高度的js代码
    js_height = "return document.body.clientHeight"
    # js_height = "return %s" % height
    try:
        driver.get(url)
        k = 1
        height = driver.execute_script(js_height)
        while True:
            if k*500 < height:
                js_move = "window.scrollTo(0,{})".format(k * 500)
                driver.execute_script(js_move)
                time.sleep(0.2)
                height = driver.execute_script(js_height)
                k += 1
            else:
                break
        driver.save_screenshot(outfile)
        print("save screenshot to {} success".format(outfile))
        time.sleep(0.1)
    except Exception as e:
        print(outfile,e)