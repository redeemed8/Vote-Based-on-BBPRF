const app = getApp()

Page({
  data: {
    avatarUrl: '',
    nickname: '',
    qnum: 0,
    vnum: 0
  },

  onLoad(options) {
  },

  onShow() {
    var that = this
    if (!app.globalData.isLogin) {
      wx.navigateTo({
        url: '/pages/unlogin/unlogin'
      })
      wx.showToast({
        title: '请先登录',
        icon: 'none',
        duration: 2000
      });
      setTimeout(function () {
        wx.hideToast();
      }, 2000);
      return
    }

    this.setData({
      avatarUrl: app.globalData.avatar_url,
      nickname: app.globalData.nickname
    });
    // ------------------
    wx.request({
      url: app.globalData.request_header + '/user/get/publish/num',
      method: 'GET',
      header: {
        "X-auth": app.globalData.LoginToken
      },
      success(res) {
        if (res.data.code !== 200) {
          wx.showToast({
            title: res.data.msg,
            icon: 'none',
            duration: 1500
          })
          return
        }
        //  ------
        that.setData({
          qnum: res.data.data.q_num,
          vnum: res.data.data.v_num
        });
      },
      fail(err) {

      }
    });

  },

  onChooseAvatar(e) {
    const { avatarUrl } = e.detail;
    app.globalData.avatar_url = avatarUrl;
    this.setData({
      avatarUrl: app.globalData.avatar_url,
    });
    //  将头像url保存到本地缓存
    wx.setStorageSync('avatarUrl', avatarUrl);
  },

  exit(e) {
    var that = this
    wx.showModal({
      title: '退出登录',
      content: '确定要退出吗？',
      confirmText: '退出',
      cancelText: '取消',
      success(res) {
        if (res.confirm) {
          that.setData({
            buttonStyle: that.data.show,
            textStyle: that.data.unshow
          })
          app.globalData.isLogin = false;
          app.globalData.LoginToken = '';
          app.globalData.question_reflash = true
          app.globalData.vote_reflash = true
          app.globalData.vote_detail_refalsh = true
          app.globalData.votePublish = false
          app.globalData.votePublishIdx = -1
          app.globalData.delVoteFlag = false
          app.globalData.delVoteIdx = -1
          app.globalData.loginToSW = ''
          app.globalData.uc = '';
          app.globalData.uc_sign = '';
          app.globalData.avatar_url = '';
          app.globalData.nickname = '';

          wx.navigateTo({
            url: '/pages/unlogin/unlogin',
          });
        }
      }
    });
  },

  set(e) {
    wx.navigateTo({
      url: '/pages/personal/personal',
    });
  },

  pro(e) {
    wx.showToast({
      title: '暂未开放',
      icon: 'none',
      duration: 1500
    });
  },

  sug(e) {
    wx.showToast({
      title: '暂未开放',
      icon: 'none',
      duration: 1500
    });
  }

})