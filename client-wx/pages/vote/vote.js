const app = getApp()

Page({
  data: {
    keyword: '',
    loading: false,
    voteList: [],
    isShow: '',
    showLoading: false
  },

  newVote(e) {
    wx.navigateTo({
      url: '/pages/new_vote/new_vote',
    })
  },

  onKeywordSearch(e) {
    this.setData({
      keyword: e.detail.value
    })
  },

  searchVote(e) {
    this.setData({
      loading: true
    })
    var that = this
    wx.showLoading({
      title: '搜索中',
      mask: true
    })
    wx.request({
      url: app.globalData.request_header + '/vote/search',
      method: 'GET',
      header: {
        "X-auth": app.globalData.LoginToken
      },
      data: {
        key: that.data.keyword
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
        that.setData({
          voteList: res.data.data
        })
      },
      fail(err) {
        wx.showToast({
          title: '服务异常，请稍后再试',
          icon: 'none',
          duration: 1000
        })
      },
      complete() {
        that.setData({
          loading: false
        })
        wx.hideLoading()
      }
    })
  },

  voteDetail(e) {
    let vid = e.currentTarget.dataset.id;
    let idx = e.currentTarget.dataset.idx;
    wx.navigateTo({
      url: '/pages/detail_vote/detail_vote?vid=' + vid + '&idx=' + idx
    })
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
    var that = this
    if (!app.globalData.isLogin) {
      that.setData({
        isShow: 'display: none;'
      });
      wx.switchTab({
        url: '/pages/index/index'
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

    // ----------------------------------
    wx.showLoading({
      title: '加载中',
      mask: true
    });
    if (!this.data.showLoading && app.globalData.vote_reflash) {
      this.setData({ showLoading: true });
      app.globalData.vote_reflash = !app.globalData.vote_reflash;
      wx.request({
        url: app.globalData.request_header + '/vote/search',
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
          that.setData({
            voteList: res.data.data
          })
        },
        fail(err) {
          wx.showToast({
            title: '服务异常，请稍后再试',
            icon: 'none',
            duration: 1000
          })
          app.globalData.vote_reflash = !app.globalData.vote_reflash
        },
        complete() {
          that.setData({ showLoading: false });
        }
      })
    }
    //  ------------------
    if (app.globalData.votePublish) {
      app.globalData.votePublish = false;
      let idx = app.globalData.votePublishIdx;
      let list = this.data.voteList;
      if (list[idx].status === 0) {
        list[idx].status = 1;
      } else {
        list[idx].status = 0;
      }
      this.setData({ voteList: list });
      app.globalData.votePublishIdx = -1;
    }
    if (app.globalData.delVoteFlag) {
      app.globalData.delVoteFlag = false;
      let list = this.data.voteList;
      let index = app.globalData.delVoteIdx;
      list.splice(index, 1);
      this.setData({ voteList: list });
      app.globalData.delVoteIdx = -1;
    }

    wx.hideLoading();
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
    app.globalData.vote_reflash = true;
    wx.switchTab({
      url: '/pages/vote/vote',
      success() {
        const currentPage = getCurrentPages().pop();
        if (currentPage) {
          currentPage.onShow();
        }
      }
    });
    wx.stopPullDownRefresh();
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