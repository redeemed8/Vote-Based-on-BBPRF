const app = getApp()

Page({
  data: {
    isShow: '',
    title: '',
    options: ['', ''],
    unshow: 'display: none;',
    hidden: 'visibility: hidden;',
    high_image_style: 'display: none;',
    title_image_style: 'display: none;',
    is_mutli_choice: false,
    agreeChecked: false,
    loading: false
  },

  radioChange(e) {
    this.setData({
      agreeChecked: !this.data.agreeChecked,
      is_mutli_choice: !this.data.is_mutli_choice
    })
  },

  titleChange(e) {
    this.setData({
      title: e.detail.value
    });
  },

  optionChange(e) {
    const index = e.currentTarget.dataset.index;
    const value = e.detail.value;
    const options = this.data.options;
    options[index] = value;
    this.setData({
      options: options
    });
  },

  addOption(e) {
    const options = this.data.options;
    if (options.length >= 10) {
      wx.showToast({
        title: '最多只能有10个选项',
        icon: 'none',
        duration: 1500
      });
      return
    }
    options.push('');
    this.setData({
      options: options
    });

    if (options.length > 1) {
      this.setData({
        title_image_style: this.data.hidden,
        high_image_style: ''
      });
    } else {
      this.setData({
        title_image_style: this.data.unshow,
        high_image_style: this.data.unshow
      });
    }
  },

  newOption(e) {
    if (this.data.title === '') {
      wx.showToast({
        title: '标题不能为空',
        icon: 'none',
        duration: 1500
      });
      return
    }
    const options = this.data.options;
    var effective_options = [];
    for (var i = 0; i < options.length; i++) {
      if (options[i] !== '') {
        effective_options.push(options[i]);
      }
    }
    if (effective_options.length < 2) {
      wx.showToast({
        title: '至少应有两个非空选项',
        icon: 'none',
        duration: 1500
      });
      return
    }

    var that = this;
    var imc = 0;
    if (this.data.is_mutli_choice) {
      imc = 1;
    }

    if (this.data.loading) {
      return
    }

    this.setData({
      loading: true
    })
    wx.showLoading({
      title: '投票创建中',
      mask: true
    })

    wx.request({
      url: app.globalData.request_header + '/vote/new',
      method: 'POST',
      header: {
        "X-auth": app.globalData.LoginToken
      },
      data: {
        "vote_title": that.data.title,
        "options": effective_options,
        "is_multi_choice": imc
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
        app.globalData.vote_reflash = true;
        wx.switchTab({
          url: '/pages/vote/vote',
        });
        setTimeout(function () {
          wx.showToast({
            title: '创建成功',
            icon: 'none',
            duration: 1000
          });
        }, 600);
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
        })
      }
    });
  },

  delOptionBtn(e) {
    if (this.data.options.length <= 2) {
      return
    }
    const index = e.currentTarget.dataset.index;
    const options = this.data.options;
    options.splice(index, 1);
    this.setData({
      options: options
    });
    if (options.length <= 2) {
      this.setData({
        title_image_style: this.data.unshow,
        high_image_style: this.data.unshow
      });
    }
    return
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