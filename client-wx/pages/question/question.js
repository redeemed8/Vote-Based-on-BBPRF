const app = getApp()

Page({
  data: {
    keyword: '',
    loading: false,
    questionList: [],
    isShow: '',
    showLoading: false
  },

  newQuestion(e) {
    wx.navigateTo({
      url: '/pages/new_question/new_question',
    })
  },

  onKeywordSearch(e) {
    this.setData({
      keyword: e.detail.value
    })
  },

  searchQuestion(e) {
    if (this.data.loading) {
      return
    }
    this.setData({
      loading: true
    })
    var that = this
    wx.showLoading({
      title: '搜索中',
      mask: true
    })
    wx.request({
      url: app.globalData.request_header + '/question/search',
      method: 'GET',
      header: {
        "X-auth": app.globalData.LoginToken
      },
      data: {
        key: that.data.keyword
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
        that.setData({
          questionList: res.data.data
        })
      },
      fail(err) {
        wx.hideLoading()
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
      }
    })
  },

  questionDetail(e) {
    var qid = e.currentTarget.dataset.id;
    var title = e.currentTarget.dataset.title;
    wx.navigateTo({
      url: '/pages/detail_questionnaire/detail_questionnaire?qid=' + qid + '&title=' + title
    })
  },

  /**
   * 生命周期函数--监听页面加载
   */
  onLoad: function (options) {

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
    if (!this.data.showLoading && app.globalData.question_reflash) {
      this.setData({ showLoading: true });
      app.globalData.question_reflash = !app.globalData.question_reflash
      wx.request({
        url: app.globalData.request_header + '/question/search',
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
            questionList: res.data.data
          })
        },
        fail(err) {
          wx.showToast({
            title: '服务异常，请稍后再试',
            icon: 'none',
            duration: 1000
          })
          app.globalData.question_reflash = !app.globalData.question_reflash
        },
        complete() {
          that.setData({ showLoading: false });
        }
      })
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
    app.globalData.question_reflash = true;
    wx.switchTab({
      url: '/pages/question/question',
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