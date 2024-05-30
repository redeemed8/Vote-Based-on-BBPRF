const app = getApp();

Page({
  data: {
    nickname: '',
    loading: false
  },

  changeNickname(e) {
    this.setData({
      nickname: e.detail.value
    });

  },

  update() {
    if (this.data.loading) {
      return;
    }
    this.setData({ loading: true });
    wx.showLoading({
      title: '修改中',
      mask: true
    });
    var that = this;
    wx.request({
      url: app.globalData.request_header + '/user/update/name/avatar',
      method: 'POST',
      data: {
        "avatar_url": '',
        "nickname": that.data.nickname
      },
      header: {
        "X-auth": app.globalData.LoginToken
      },
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
        // -----
        app.globalData.nickname = that.data.nickname;
        wx.switchTab({
          url: '/pages/my/my',
        });
        wx.showToast({
          title: res.data.data,
          icon: 'none',
          duration: 1500
        });
      },
      fail(err) {
        wx.hideLoading();
        wx.showToast({
          title: '服务异常，请稍后再试',
          icon: 'none',
          duration: 1000
        });
      },
      complete() {
        that.setData({ loading: false });
      }
    });
  },

  onLoad(options) {

  },

  onShow() {

  }

});