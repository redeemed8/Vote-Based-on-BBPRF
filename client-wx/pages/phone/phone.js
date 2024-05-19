const app = getApp()

Page({
  data: {
    phoneNumber: '',
    verifyCode: '',
    disableButton: false,
    buttonText: '获取验证码',
    loginLoading: false,
    toPath: ''
  },

  inputPhone(e) {
    this.setData({
      phoneNumber: e.detail.value
    })
  },

  inputVerify(e) {
    this.setData({
      verifyCode: e.detail.value
    })
  },

  getVerifyCode() {
    if (!this.data.disableButton) {
      this.setData({
        disableButton: true,
        buttonText: '60s可点击'
      })
      var that = this
      wx.showLoading({
        title: '发送中',
        mask: true
      })
      wx.request({
        url: app.globalData.request_header + '/user/send/code',
        method: 'GET',
        success(res) {
          wx.hideLoading();
          if (res.data.code !== 200) {
            wx.showToast({
              title: res.data.msg,
              icon: 'none',
              duration: 1500
            })
            return
          }
          var countDownSeconds = 60;
          var timer = setInterval(() => {
            countDownSeconds--;
            if (countDownSeconds > 0) {
              that.setData({
                buttonText: countDownSeconds + 's可点击'
              });
            } else {
              clearInterval(timer);
              that.setData({
                disableButton: false,
                buttonText: '获取验证码',
              });
            }
          }, 1000);
          if (res.data.code === 200) {
            wx.showToast({
              title: '验证码已发送',
              icon: 'none',
              duration: 1000
            })
          }
        },
        fail(err) {
          wx.hideLoading();
          wx.showToast({
            title: '验证码发送失败',
            icon: 'none',
            duration: 1000
          })
        }
      });
    }
  },

  phoneLogin(e) {
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
        url: app.globalData.request_header + '/user/login/phone',
        method: 'POST',
        data: {
          phone: that.data.phoneNumber,
          code: that.data.verifyCode
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