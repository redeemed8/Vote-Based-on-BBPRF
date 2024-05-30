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
    loading: false,
    imageUrl: ''
  },

  chooseImage: function () {
    wx.chooseImage({
      count: 1,
      sizeType: ['original', 'compressed'],
      sourceType: ['album', 'camera'],
      success: (res) => {
        const tempFilePath = res.tempFilePaths[0];
        this.setData({
          imageUrl: tempFilePath
        });
      }
    });
  },

  uploadImage: function (filePath, title, options, imc) {
    wx.uploadFile({
      url: app.globalData.request_header + '/vote/new',
      filePath: filePath,
      name: 'file',
      header: {
        "X-auth": app.globalData.LoginToken
      },
      formData: {
        "vote_title": title,
        "options": JSON.stringify(options),
        "is_multi_choice": String(imc)
      },
      success: (res) => {
        wx.hideLoading();
        var dataObj = JSON.parse(res.data);
        if (dataObj.code !== 200) {
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
      fail: (err) => {
        wx.hideLoading();
        wx.showToast({
          title: '服务异常，请稍后再试',
          icon: 'none',
          duration: 1000
        });
        return
      }
    });
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
    var that = this
    wx.showModal({
      title: '是否开启多选',
      confirmText: '开启',
      cancelText: '不开启',
      success(res) {
        if (res.confirm) {
          that.setData({ is_mutli_choice: true });
          that.doNewOption();
        } else {
          that.setData({ is_mutli_choice: false });
          that.doNewOption();
        }
      }
    });
  },

  doNewOption() {
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

    var imc = 0;
    if (this.data.is_mutli_choice) {
      imc = 1;
    }

    if (this.data.loading) {
      return
    }

    if (this.data.imageUrl === '') {
      wx.showToast({
        title: '请选择一张图片',
        icon: 'none',
        duration: 1500
      });
      return
    }

    this.setData({
      loading: true
    })
    wx.showLoading({
      title: '投票创建中',
      mask: true
    })
    this.uploadImage(this.data.imageUrl, this.data.title, effective_options, imc);
    this.setData({ loading: false });
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

  onLoad(options) {

  },

  onShow() {
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
  }

});