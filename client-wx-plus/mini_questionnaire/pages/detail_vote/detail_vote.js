import * as echarts from '../../ec-canvas/echarts';
const app = getApp()

Page({
  data: {
    vid: '',
    isShow: '',
    loading: false,
    optionNames: [],
    optionCounts: [],
    title: '',
    status: 0,
    participants: 0,
    imc: 0,
    optionL: ['A', 'B', 'C', 'D', 'E', 'F', 'G', 'H', 'I', 'J'],
    btnText: '显示统计图',
    lazyEc: {
      lazyLoad: true
    },
    ecLoaded: false,
    ecShow: '',
    rotate: 0,
    publishBtnText: '发布投票',
    tj_color: 'color: rgba(128, 128, 128, 1);',
    index: -1,
    url: '',
  },

  showCount(e) {
    if (!this.data.ecLoaded) {
      this.setData({
        ecLoaded: true,
        btnText: '统计图',
        tj_color: 'color: rgba(208, 109, 90, 1);'
      });
      this.init()
      return
    }
    if (this.data.btnText === '显示统计图') {
      this.setData({
        btnText: '统计图',
        ecShow: '',
        tj_color: 'color: rgba(208, 109, 90, 1);'
      });
    } else {
      this.setData({
        btnText: '显示统计图',
        ecShow: 'display: none;',
        tj_color: 'color: rgba(128, 128, 128, 1);'
      });
    }
  },

  onLoad(options) {
    app.globalData.vote_detail_refalsh = true;
    this.setData({
      vid: options['vid'],
      index: options['idx']
    });
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

    //  ----------------------------------------
    wx.showLoading({
      title: '加载中',
      mask: true
    });
    if (!this.data.loading && app.globalData.vote_detail_refalsh) {
      this.setData({ showLoading: true });
      app.globalData.vote_detail_refalsh = !app.globalData.vote_detail_refalsh;
      wx.request({
        url: app.globalData.request_header + '/vote/get-detail',
        method: 'GET',
        header: {
          "X-auth": app.globalData.LoginToken
        },
        data: {
          "id": that.data.vid
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
          if (res.data.data.ans_count === null || res.data.data.ans_count === undefined) {
            return
          }
          var counts = []
          for (var i = 0; i < res.data.data.ans_count.length; i++) {
            counts.push(res.data.data.ans_count[i].option_count);
          }

          console.log(res.data);
          that.setData({
            optionNames: res.data.data.options,
            optionCounts: counts,
            title: res.data.data.title,
            status: res.data.data.status,
            participants: res.data.data.participants,
            imc: res.data.data.is_multi_choice,
            url: res.data.data.url,
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
          that.setData({ showLoading: false });
        }
      });
    }
    // ------------------------
    // 获取到组件
    this.lazyComponent = this.selectComponent('#lazy-mychart-dom')
    wx.hideLoading();
  },
  init() {
    this.lazyComponent.init((canvas, width, height, dpr) => {
      let chart = echarts.init(canvas, null, {
        width: width,
        height: height,
        devicePixelRatio: dpr
      })
      let option = getOption()
      chart.setOption(option)
      this.chart = chart
      return chart
    })
  },

  radioChange(e) {
    if (this.data.imc === 0) {
      this.setData({ imc: 1 });
    } else {
      this.setData({ imc: 0 });
    }
  },

  publish(e) {
    if (this.data.vid === '' || this.data.loading) {
      wx.showToast({
        title: '操作太频繁，请稍后再试',
        icon: 'none',
        duration: 1000
      });
      return;
    }
    var that = this;
    this.setData({ loading: true });
    wx.showLoading({ title: '加载中', mask: true });
    wx.request({
      url: app.globalData.request_header + '/vote/update/status',
      method: 'GET',
      header: { "X-auth": app.globalData.LoginToken },
      data: { "vid": that.data.vid, "imc": that.data.imc },
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
          status: res.data.data.status,
          imc: res.data.data.imc
        });
        app.globalData.votePublish = !app.globalData.votePublish;
        app.globalData.votePublishIdx = that.data.index;

        if (that.data.status === 0) {
          wx.showToast({
            title: '已关闭发布',
            icon: 'none',
            duration: 1500
          })
        } else {
          wx.showToast({
            title: '已开启发布',
            icon: 'none',
            duration: 1500
          })
        }
      },
      fail(err) {
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
    wx.hideLoading();
  },

  delete(e) {
    var that = this;
    wx.showModal({
      title: '删除投票',
      content: '确定要删除吗？',
      confirmText: '删除',
      cancelText: '取消',
      success(res) {
        if (res.confirm) {
          that.deleteVote();
        }
      }
    });
  },

  deleteVote(e) {
    if (this.data.vid === '' || this.data.loading) {
      wx.showToast({
        title: '操作太频繁，请稍后再试',
        icon: 'none',
        duration: 1000
      });
      return;
    }
    var that = this;
    this.setData({ loading: true });
    wx.showLoading({ title: '加载中', mask: true });
    wx.request({
      url: app.globalData.request_header + '/vote/del',
      method: 'GET',
      header: { "X-auth": app.globalData.LoginToken },
      data: { "id": that.data.vid },
      success(res) {
        if (res.data.code !== 200) {
          wx.showToast({
            title: res.data.msg,
            icon: 'none',
            duration: 1500
          })
          return
        }
        app.globalData.delVoteFlag = true;
        app.globalData.delVoteIdx = that.data.index;
        wx.showToast({
          title: '删除成功',
          icon: 'none',
          duration: 1000
        });
        setTimeout(function () {
          wx.switchTab({
            url: '/pages/vote/vote',
          })
        }, 500);
      },
      fail(err) {
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
    wx.hideLoading();
  },

  //  分享
  onShareAppMessage() {
    return {
      title: '投票',
      path: '/pages/voting/voting?vid=' + this.data.vid,
      imageUrl: '',
      desc: '快来参与投票吧!'
    }
  }
});

function getOption() {
  const pages = getCurrentPages();
  const currentPage = pages[pages.length - 1];
  const data = currentPage.data;

  let names = []
  for (let i = 0; i < data.optionNames.length; i++) {
    names.push(truncateString(data.optionNames[i], 6));
  }

  if (data.optionNames.length > 5) {
    data.rotate = 25;
  } else {
    data.rotate = 0;
  }

  return {
    xAxis: {
      type: 'category',
      data: names,
      axisLine: {
        lineStyle: {
          width: 17, // 设置坐标轴线的宽度
          color: '#4b4b4b' // 设置坐标轴线的颜色
        }
      },
      axisLabel: {
        fontSize: 12,
        color: '',
        interval: 0,
        margin: 30,
        rotate: data.rotate,
        align: 'center' // 居中对齐
      }
    },
    yAxis: {
      type: 'value',
      axisLabel: {
        margin: 15
      },
      axisLine: {
        lineStyle: {
          color: '#4b4b4b', // 设置坐标轴线的颜色
        }
      }
    },
    series: [{
      name: '星期',
      type: 'bar',
      itemStyle: {
        normal: {
          barBorderRadius: [6, 6, 0, 0],
          color: new echarts.graphic.LinearGradient(
            0, 0, 0, 1,
            [
              { offset: 0, color: '#ffd700' },
              { offset: 0.5, color: '#ffa07a' },
              { offset: 1, color: '#d06d5a' }
            ]
          ),
        }
      },
      label: {
        show: true,
        position: 'top',
        fontSize: 11,
        formatter: function (params) {
          return params.value;
        }
      },
      emphasis: {
        focus: 'series'
      },
      data: data.optionCounts
    }]
  }

}

function truncateString(str, maxLength) {
  let length = 0;
  let truncated = '';

  for (let i = 0; i < str.length; i++) {
    length += (str.charCodeAt(i) > 255) ? 2 : 1;
    if (length >= maxLength) {
      truncated = str.slice(0, i + 1);
      if (i !== str.length - 1) {
        truncated += "..";
      }
      break;
    }
    truncated += str[i];
  }
  return truncated || str;
}
