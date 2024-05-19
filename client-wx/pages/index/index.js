const app = getApp()

Page({
  data: {
    request_prefix: 'http://localhost:3656',
    unshow: 'display:none;',
    show: 'display:block;',
    buttonStyle: 'display:none;',
    textStyle: 'display:none;',
    toPath: ''
  },

  //  账号登录
  accountLogin() {
    let url = '/pages/account/account';
    if (this.data.toPath !== '' && this.data.toPath !== undefined) {
      url += '?topath=' + this.data.toPath;
    }

    wx.navigateTo({
      url: url,
    })
  },

  //  手机号登录
  phoneLogin() {
    let url = '/pages/phone/phone';
    if (this.data.toPath !== '' && this.data.toPath !== undefined) {
      url += '?topath=' + this.data.toPath;
    }

    wx.navigateTo({
      url: url,
    })
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

          wx.reLaunch({
            url: '/pages/index/index'
          });
        }
      }
    });
  },

  /**
   * 生命周期函数--监听页面加载
   */
  onLoad(options) {
    let topath = app.globalData.loginToSW;
    if (topath === null || topath === undefined) {
      return;
    }
    this.setData({ toPath: topath }); //  保存跳转路径

    load_pk_s();

  },

  /**
   * 生命周期函数--监听页面显示
   */
  onShow() {
    if (app.globalData.isLogin) {
      this.setData({
        buttonStyle: this.data.unshow,
        textStyle: 'display:flex;'
      })
    } else {
      this.setData({
        buttonStyle: this.data.show,
        textStyle: this.data.unshow
      })
    }
  }
});


function load_pk_s() {
  wx.showLoading({
    title: '载入sss中',
    mask: true
  });
  if (app.globalData.pk_s.N > 0) {
    wx.hideLoading();
    return
  }
  //  如果未加载就加载
  wx.request({
    url: app.globalData.request_header + '/auth/get-pprms',
    method: 'GET',
    success(res) {
      app.globalData.pk_s.N = res.data.data.N
      app.globalData.pk_s.G = res.data.data.G
      app.globalData.pk_s.Y = res.data.data.Y
      app.globalData.pk_s.H = res.data.data.H
      app.globalData.pk_s.ct_u = res.data.data.CtU
      app.globalData.pk_s.ct_y = res.data.data.CtY
    },
    fail(err) {
      wx.hideLoading();
      wx.showToast({
        title: '小程序初始化失败!',
        icon: 'none',
        duration: 1000
      });
    }
  });
  wx.hideLoading();
}

