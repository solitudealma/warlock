class Settings {
    constructor(root) {
        this.root = root;
        this.platform = "WEB";
        if(this.root.AcWingOS) this.platform = "ACAPP";
        this.username = "";
        this.photo = "";
        this.$settings = $(`
            <div class="ac-game-settings">
                <div class="ac-game-settings-login">
                    <div class="ac-game-settings-title">
                        登录
                    </div>
                    <div class="ac-game-settings-username">
                        <div class="ac-game-settings-item">
                            <input type="text" placeholder="用户名">
                        </div>
                    </div>
                    <div class="ac-game-settings-password">
                        <div class="ac-game-settings-item">
                            <input type="password" placeholder="密码">
                        </div>
                    </div>
                    <div class="ac-game-settings-submit">
                        <div class="ac-game-settings-item">
                            <button>登录</button>
                        </div>
                    </div>
                    <div class="ac-game-settings-error-message">
                    </div>
                    <div class="ac-game-settings-option">
                        注册
                    </div>
                    <!-- display: inline; 会影响下一行 -->
                    <br>
<!--                    <div class="ac-game-settings-acwing">-->
<!--                        <img width="30" src="https://acapp.yt1209.com:4430/static/images/settings/acwing_logo.png" -->
<!--                             alt="">-->
<!--                        <br>-->
<!--                        <div>-->
<!--                            AcWing一键登录-->
<!--                        </div>-->
<!--                    </div>-->
                </div>
                <div class="ac-game-settings-register">
                    <div class="ac-game-settings-title">
                        注册
                    </div>
                    <div class="ac-game-settings-username">
                        <div class="ac-game-settings-item">
                            <input type="text" placeholder="用户名">
                        </div>
                    </div>
                    <div class="ac-game-settings-password ac-game-settings-password-first">
                        <div class="ac-game-settings-item">
                            <input type="password" placeholder="密码">
                        </div>
                    </div>
                    <div class="ac-game-settings-password ac-game-settings-password-second">
                        <div class="ac-game-settings-item">
                            <input type="password" placeholder="确认密码">
                        </div>
                    </div>
                    <div class="ac-game-settings-submit">
                        <div class="ac-game-settings-item">
                            <button>注册</button>
                        </div>
                    </div>
                    <div class="ac-game-settings-error-message">
                    </div>
                    <div class="ac-game-settings-option">
                        登录
                    </div>
                    <!-- display: inline; 会影响下一行 -->
                    <br>
<!--                    <div class="ac-game-settings-acwing">-->
<!--                        <img width="30" src="https://acapp.yt1209.com:4430/static/images/settings/acwing_logo.png" -->
<!--                             alt="">-->
<!--                        <br>-->
<!--                        <div>-->
<!--                            AcWing一键登录-->
<!--                        </div>-->
<!--                    </div>-->
                </div>
            </div>
            `);
        this.$login = this.$settings.find(".ac-game-settings-login");
        this.$login_username = this.$login.find(".ac-game-settings-username input");
        this.$login_password = this.$login.find(".ac-game-settings-password input");
        this.$login_submit = this.$login.find(".ac-game-settings-submit button");
        this.$login_error_message = this.$login.find(".ac-game-settings-error-message");
        this.$login_register = this.$login.find(".ac-game-settings-option");

        this.$login.hide();

        this.$register = this.$settings.find(".ac-game-settings-register");
        this.$register_username = this.$register.find(".ac-game-settings-username input");
        this.$register_password = this.$register.find(".ac-game-settings-password-first input");
        this.$register_password_confirm = this.$register.find(".ac-game-settings-password-second input");
        this.$register_submit = this.$register.find(".ac-game-settings-submit button");
        this.$register_error_message = this.$register.find(".ac-game-settings-error-message");
        this.$register_login = this.$register.find(".ac-game-settings-option");

        this.$register.hide();

        this.$acwing_login = this.$settings.find(".ac-game-settings-acwing img")
        this.root.$ac_game.append(this.$settings);
        this.start();
    }

    start() {
        if(this.platform === "WEB") {
            this.getinfo_web();
            this.add_listening_events();
        } else {
            this.getinfo_acapp();
        }
    }

    add_listening_events() {
        let outer = this;
        this.add_listening_events_login();
        this.add_listening_events_register();

        this.$acwing_login.click(function() {
            outer.acwing_login();
        })
    }

    add_listening_events_login() {
        let outer = this;
        this.$login_register.click(function() {
            outer.register();
        });
        this.$login_submit.click(function() {
            outer.login_on_remote();
        });
    }

    add_listening_events_register() {
        let outer = this;
        this.$register_login.click(function() {
            outer.login();
        });
        this.$register_submit.click(function() {
            outer.register_on_remote();
        });
    }

    acwing_login() {
        $.ajax({
            url: "https://acapp.yt1209.com:4430/settings/acwing/web/apply-code",
            method: "GET",
            data: {},
            success: function(res) {
                if(res.result === "success") {
                    window.location.replace(res.apply_code_url);
                }
            }
        })
    }

    login_on_remote() { // 在远程服务器上登录
        let outer = this;
        let username = this.$login_username.val();
        let password = this.$login_password.val();
        this.$login_error_message.empty();
        $.ajax({
            url: "/settings/login",
            method: "GET",
            data: {
                username: username,
                password: password,
            },
            success: function(res) {
                if(res.msg === "登录成功") {
                    localStorage.setItem('token', res.data.token);
                    localStorage.setItem('username', res.data.user.username);
                    location.reload();
                } else {
                    localStorage.clear();
                    outer.$login_error_message.text(res.msg);
                }
            }
        })
    }

    register_on_remote() { // 在远程服务器上注册
        let outer = this;
        let username = this.$register_username.val();
        let password = this.$register_password.val();
        let password_confirm = this.$register_password_confirm.val();
        this.$register_error_message.empty();

        $.ajax({
            url: "/settings/register",
            method: "GET",
            data: {
                username: username,
                password: password,
                password_confirm: password_confirm,
            },
            success: function(res) {
                console.log(res)
                if(res.msg === "注册成功") {
                    location.reload(); // 刷新页面
                } else {
                    outer.$register_error_message.text(res.msg);
                }
            }
        })
    }

    logout_on_remote() { // 在远程服务器上登出
        if(this.platform === "ACAPP") return false;
        let token = localStorage.getItem('token');
        $.ajax({
            url: "/settings/logout",
            method: "GET",
            headers: {
              'x-token': token,
            },
            data: {},
            success: function(res) {
                if(res.msg === "jwt作废成功") {
                    localStorage.clear()
                    location.reload();
                }
            }
        })
    }

    register() { // 打开注册界面
        this.$login.hide();
        this.$register.show();
    }

    login() { // 打开登录界面
        this.$register.hide();
        this.$login.show();
    }

    acapp_login(appid, redirect_uri, scope, state) {
        let outer = this;
        this.root.AcWingOS.api.oauth2.authorize(appid, redirect_uri, scope, state, function(res) {
            if(res.result === "success") {
                outer.username = res.username;
                outer.photo = res.photo;
                outer.hide();
                outer.root.menu.show();
            }
        });
    }

    getinfo_acapp() {
        let outer = this;
        $.ajax({
            url: "https://acapp.yt1209.com:4430/settings/acwing/acapp/apply-code",
            method: "GET",
            success: function(res) {
                if(res.result === "success") {
                    outer.acapp_login(res.appid, res.redirect_uri, res.scope, res.state);
                }
            },
            error: function(req , err) {
            }
        });
    }

    getinfo_web() {
        let outer = this;
        let token = localStorage.getItem('token');
        if (typeof token != null) {
            $.ajax({
                url: "/settings/getinfo",
                method: "GET",
                headers: {
                    // 如果后台没有跨域处理，这个自定义
                    "x-token": token,
                },
                data: {
                    platform: this.platform,
                },
                success: function(res) {
                    if(res.msg === "success") {
                        outer.username = res.data.userInfo.username;
                        outer.photo = res.data.userInfo.avatar;
                        outer.hide();
                        outer.root.menu.show();
                    } else {
                        localStorage.clear()
                        outer.login();
                    }
                }
            });
        }
    }

    show() {
        this.$settings.show();
    }

    hide() {
        this.$settings.hide();
    }
}
