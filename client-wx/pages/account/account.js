const app = getApp()

Page({
  data: {
    account: '',
    password: '',
    loginLoading: false,
    toPath: ''
  },

  inputAccount(e) {
    this.setData({
      account: e.detail.value
    })
  },

  inputPassword(e) {
    this.setData({
      password: e.detail.value
    })
  },

  register(e) {
    let url = '/pages/register/register';
    if (this.data.toPath !== '' && this.data.toPath !== undefined) {
      url += '?topath=' + this.data.toPath;
    }

    wx.navigateTo({
      url: url,
    })
  },

  accountLogin(e) {
    const md5 = require('../../libs/md5.js');
    var that = this
    if (!this.data.loginLoading) {
      this.setData({
        loginLoading: true
      });
      wx.showLoading({
        title: '登录中',
        mask: true
      })
      wx.request({
        url: app.globalData.request_header + '/user/login/account',
        method: 'POST',
        data: {
          account: that.data.account,
          password: md5(that.data.password)
        },
        success(res) {
          wx.hideLoading();
          if (res.data.code === 200) {
            app.globalData.isLogin = true
            app.globalData.LoginToken = res.data.data

            LoadUcSecret();

            //  判断是否有目的路径
            if (that.data.toPath === '' || that.data.toPath === undefined) {
              wx.switchTab({
                url: '/pages/index/index',
              });
            } else {
              wx.navigateTo({
                url: decodeURIComponent(that.data.toPath),
              });
            }

            wx.showToast({
              title: '登录成功',
              icon: 'none',
              duration: 1500
            });
          } else {
            wx.showToast({
              title: res.data.msg,
              icon: 'none',
              duration: 1000
            });
          }
        },
        fail(err) {
          wx.hideLoading();
          wx.showToast({
            title: '服务异常，请稍后再试',
            icon: 'none',
            duration: 1000
          })
        },
        complete() {
          that.setData({
            loginLoading: false
          });
        }
      });
    }
  },

  /**
   * 生命周期函数--监听页面加载
   */
  onLoad(options) {
    let topath = options['topath'];
    if (topath === null || topath === undefined) {
      return;
    }
    this.setData({ toPath: topath }); //  保存跳转路径
  }
 
});

function LoadUcSecret() {
  wx.request({
    url: app.globalData.request_header2 + '/pkcs/get/signed/privk',
    method: 'POST',
    data: {
      "identify": app.globalData.LoginToken
    },
    header: {
      'Content-Type': 'application/x-www-form-urlencoded'
    },
    success(res) {
      wx.hideLoading();
      app.globalData.uc = res.data.data.priv_key
      app.globalData.uc_sign = res.data.data.sign
    },
    fail(err) {
      wx.hideLoading();
      wx.showToast({
        title: '获取客服端密钥失败!',
        icon: 'none',
        duration: 1000
      });
    }
  });
};