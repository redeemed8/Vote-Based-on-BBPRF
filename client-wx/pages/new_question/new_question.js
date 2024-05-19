const app = getApp()

Page({
  data: {
    title: '',
    titleMaxLength: 100
  },

  inputTitle(e) {
    this.setData({
      title: e.detail.value
    });
  },

  newQuestion(e) {
    if (this.data.title === '' || this.data.title.length < 1) {
      wx.showToast({
        title: '标题不能为空',
        icon: 'none',
        duration: 1000
      })
      return
    }
    var that = this;
    wx.navigateTo({
      url: '/pages/add_question/add_question?title=' + that.data.title,
    });
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
    }
    return
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