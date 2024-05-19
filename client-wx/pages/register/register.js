const app = getApp()

Page({
  data: {
    account: '',
    password: '',
    repassword: '',
    phoneNumber: '',
    verifyCode: '',
    disableButton: false,
    buttonText: '获取验证码',
    loading: false
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

  inputRepassword(e) {
    this.setData({
      repassword: e.detail.value
    })
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
          if (res.data.code === 200) {
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
            wx.hideLoading();
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

  registerAccount(e) {
    const md5 = require('../../libs/md5.js');
    var that = this;
    if (this.data.account === '') {
      wx.showToast({
        title: '账号不能为空',
        icon: 'none',
        duration: 1000
      });
      return
    }
    if (this.data.password === '') {
      wx.showToast({
        title: '密码不能为空',
        icon: 'none',
        duration: 1000
      });
      return
    }
    if (this.data.repassword !== this.data.password) {
      wx.showToast({
        title: '两次密码不一致',
        icon: 'none',
        duration: 1000
      });
      return
    }
    if (this.data.phoneNumber === '') {
      wx.showToast({
        title: '手机号不能为空',
        icon: 'none',
        duration: 1000
      });
      return
    }
    if (this.data.verifyCode === '') {
      wx.showToast({
        title: '验证码不能为空',
        icon: 'none',
        duration: 1000
      });
      return
    }
    wx.showLoading({
      title: '加载中',
      mask: true
    })
    if (!this.data.loading) {
      this.setData({
        loading: true
      });
      wx.request({
        url: app.globalData.request_header + '/user/register/account',
        method: 'POST',
        data: {
          account: that.data.account,
          password: md5(that.data.password),
          repassword: md5(that.data.repassword),
          phone: that.data.phoneNumber,
          code: that.data.verifyCode
        },
        success(res) {
          wx.hideLoading();
          if (res.data.code === 200) {
            wx.navigateTo({
              url: '/pages/account/account',
            });

            wx.showToast({
              title: '注册成功',
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
            loading: false
          });
        }
      })
    }
  },

  /**
   * 生命周期函数--监听页面加载
   */
  onLoad(options) {

  },

  /**
   * 生命周期函数--监听页面初次渲染完成
   */
  onReady() {

  },

  /**
   * 生命周期函数--监听页面显示
   */
  onShow() {

  },

  /**
   * 生命周期函数--监听页面隐藏
   */
  onHide() {

  },

  /**
   * 生命周期函数--监听页面卸载
   */
  onUnload() {

  },

  /**
   * 页面相关事件处理函数--监听用户下拉动作
   */
  onPullDownRefresh() {

  },

  /**
   * 页面上拉触底事件的处理函数
   */
  onReachBottom() {

  },

  /**
   * 用户点击右上角分享
   */
  onShareAppMessage() {

  }
})