const app = getApp()

Page({
  data: {
    isShow: '',
    loading: false
  },

  loadAvatar() {
    if (this.data.loading || app.globalData.avatar_url !== '') {
      return;
    }
    this.setData({ loading: true });
    //  加载用户昵称和头像
    var that = this;
    wx.request({
      url: app.globalData.request_header + '/user/get/name/avatar',
      method: 'GET',
      header: {
        'X-auth': app.globalData.LoginToken
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
        app.globalData.avatar_url = wx.getStorageSync('avatarUrl');
        app.globalData.nickname = res.data.data.nickname;
        if (app.globalData.avatar_url === "" || app.globalData.avatar_url === undefined) {
          app.globalData.avatar_url = 'https://pvs.81jcpd.cn/2024/01/28/6c/6c95e7dd0ad0027ef8909eb639011769.jpeg';
          return;
        }
        //  测试头像url的有效性
        wx.getImageInfo({
          src: app.globalData.avatar_url,
          fail: (res) => {
            // 图片加载失败，说明图片地址失效，使用默认头像
            app.globalData.avatar_url = 'https://pvs.81jcpd.cn/2024/01/28/6c/6c95e7dd0ad0027ef8909eb639011769.jpeg';
          }
        });
      },
      fail(err) {
        wx.showToast({
          title: '服务异常，请稍后再试',
          icon: 'none',
          duration: 1000
        })
      },
      complete() {
        that.setData({ loading: false });
      }
    });
  },

  onShow() {
    this.loadAvatar();
    var that = this
    if (!app.globalData.isLogin) {
      that.setData({
        isShow: 'display: none;'
      });
      wx.navigateTo({
        url: '/pages/unlogin/unlogin'
      })
      wx.showToast({
        title: '请先登录',
        icon: 'none',
        duration: 2000
      })
      setTimeout(function () {
        wx.hideToast();
      }, 2000);
      return
    } else {
      that.setData({
        isShow: 'display: block;'
      });
    }
  },

  newQ() {
    wx.showToast({
      title: '暂未开放',
      icon: 'none',
      duration: 1500
    });
  },

  newV() {
    wx.navigateTo({
      url: '/pages/new_vote/new_vote',
    })
  }
});



