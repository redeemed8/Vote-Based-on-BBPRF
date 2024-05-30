const app = getApp()

Page({
  data: {
    w: 0,
    h: 0,
    phoneNumber: '',
    verifyCode: '',
    disableButton: false,
    buttonText: '获取验证码',
    loginLoading: false,
    toPath: '',
    bigR: 0,
    circle_x: 0,
    circle_y: 0,
    width2: '',
    height2: '',
    small_x: 0,
    small_y: 0,
    small_r: 0
  },

  inputPhone(e) {
    this.setData({
      phoneNumber: e.detail.value
    })
  },

  inputVerify(e) {
    this.setData({
      verifyCode: e.detail.value
    })
  },

  getVerifyCode() {
    if (!this.data.disableButton) {
      this.setData({
        disableButton: true,
        buttonText: '60s可点击'
      })
      var that = this
      wx.showLoading({
        title: '发送中',
        mask: true
      })
      wx.request({
        url: app.globalData.request_header + '/user/send/code',
        method: 'GET',
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
          var countDownSeconds = 60;
          var timer = setInterval(() => {
            countDownSeconds--;
            if (countDownSeconds > 0) {
              that.setData({
                buttonText: countDownSeconds + 's可点击'
              });
            } else {
              clearInterval(timer);
              that.setData({
                disableButton: false,
                buttonText: '获取验证码',
              });
            }
          }, 1000);
          if (res.data.code === 200) {
            wx.showToast({
              title: '验证码已发送',
              icon: 'none',
              duration: 1000
            })
          }
        },
        fail(err) {
          wx.hideLoading();
          wx.showToast({
            title: '验证码发送失败',
            icon: 'none',
            duration: 1000
          })
        }
      });
    }
  },

  phoneLogin(e) {
    var that = this
    if (!this.data.loginLoading) {
      this.setData({
        loginLoading: true
      });
      wx.showLoading({
        title: '登录中',
        mask: true
      })
      wx.request({
        url: app.globalData.request_header + '/user/login/phone',
        method: 'POST',
        data: {
          phone: that.data.phoneNumber,
          code: that.data.verifyCode
        },
        success(res) {
          wx.hideLoading();
          if (res.data.code === 200) {
            app.globalData.isLogin = true
            app.globalData.LoginToken = res.data.data

            LoadUcSecret();

            //  判断是否有目的路径
            if (that.data.toPath === '' || that.data.toPath === undefined) {
              wx.switchTab({
                url: '/pages/unlogin/unlogin',
              });
            } else {
              wx.navigateTo({
                url: decodeURIComponent(that.data.toPath),
              });
            }

            wx.showToast({
              title: '登录成功',
              icon: 'none',
              duration: 1500
            });
          } else {
            wx.showToast({
              title: res.data.msg,
              icon: 'none',
              duration: 1000
            });
          }
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
            loginLoading: false
          });
        }
      });
    }
  },

  /**
   * 生命周期函数--监听页面加载
   */
  onLoad(options) {
    let topath = options['topath'];
    if (topath === null || topath === undefined) {
      return;
    }
    this.setData({ toPath: topath }); //  保存跳转路径
  },

  onReady() {
    var that = this;
    wx.getSystemInfo({
      success: function (res) {
        that.setData({ w: res.windowWidth });
        that.setData({ h: res.windowHeight });
        that.drawBigCircle();
        that.drawTwo();
      }
    });
  },

  drawBigCircle: function () {
    const ctx = wx.createCanvasContext('myCanvas');
    //  圆形参数
    const centerX = this.data.w / 2; //   圆心x
    const radius = this.data.w * 0.4113; //  半径

    // const radius = radius_t < this.data.h * 0.8 ? radius_t : this.data.h * 0.8;

    //  矩形参数
    const rectWidth = radius * 2;
    const rectHeight = 250;
    const rectX = centerX - radius; //  矩形高

    const centerY = radius + 2;  //  圆心y
    const rectY = centerY;

    this.setData({ bigR: radius, circle_x: centerX, circle_y: centerY });

    //  绘制半圆
    ctx.beginPath();
    ctx.arc(centerX, centerY, radius, Math.PI, 2 * Math.PI, false);
    // 绘制矩形
    ctx.rect(rectX, rectY, rectWidth, rectHeight);
    // 设置阴影颜色和模糊度
    ctx.shadowColor = 'rgba(0, 0, 0, 0.5)';
    ctx.shadowBlur = 5;
    //  绘制
    ctx.closePath();
    ctx.fillStyle = '#89e2ad';
    ctx.fill();
    ctx.draw();
  },

  drawTwo() {
    const context = wx.createCanvasContext('myCanvas2');
    const x = this.data.circle_x; // 圆心 x 坐标
    const y = this.data.circle_y - this.data.bigR / 2 + 11; // 圆心 y 坐标
    const radius = this.data.bigR / 2.7; // 圆的半径
    const shadowRadius = radius * 1.06; // 阴影圆的半径，比原圆稍大一点

    this.setData({ small_x: x, small_y: y, small_r: radius });

    // 绘制阴影圆
    context.beginPath();
    context.arc(x, y + 4, shadowRadius - 2, 0, 2 * Math.PI); // 圆心坐标 (x, y + 10)，半径 shadowRadius，起始角度 0，终止角度 2 * Math.PI
    context.setFillStyle('rgba(0, 0, 0, 0.2)'); // 设置填充颜色为灰色半透明
    context.fill(); // 填充阴影圆

    // 绘制实心圆
    context.beginPath();
    context.arc(x, y, radius, 0, 2 * Math.PI); // 圆心坐标 (x, y)，半径 radius，起始角度 0，终止角度 2 * Math.PI
    context.setFillStyle('#87eb93'); // 设置填充颜色为黑色
    context.fill(); // 填充圆

    context.setFillStyle('white'); // 设置文字颜色为白色
    context.setFontSize(26); // 设置文字大小
    context.setTextAlign('center'); // 文字居中对齐
    context.setTextBaseline('middle'); // 文字垂直居中对齐
    context.fillText('进入', x, y); // 在圆心位置绘制文字 Hello

    context.draw(); // 渲染到 canvas
  },

  // 处理点击事件
  onCanvasClick: function (e) {
    const x = e.detail.x;
    const y = e.detail.y - this.data.h * 0.743;
    const centerX = this.data.small_x;
    const centerY = this.data.small_y;
    const radius = this.data.small_r;
    // 判断点击是否在半圆内
    const rrr = Math.sqrt((x - centerX) ** 2 + (y - centerY) ** 2);
    if (rrr <= radius) {
      this.phoneLogin();
    }
  }

});


function LoadUcSecret() {
  wx.request({
    url: app.globalData.request_header2 + '/pkcs/get/signed/privk',
    method: 'POST',
    data: {
      "identify": app.globalData.LoginToken
    },
    header: {
      'Content-Type': 'application/x-www-form-urlencoded'
    },
    success(res) {
      wx.hideLoading();
      app.globalData.uc = res.data.data.priv_key
      app.globalData.uc_sign = res.data.data.sign
    },
    fail(err) {
      wx.hideLoading();
      wx.showToast({
        title: '获取客服端密钥失败!',
        icon: 'none',
        duration: 1000
      });
    }
  });
};