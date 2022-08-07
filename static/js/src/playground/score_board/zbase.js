class ScoreBoard extends AcGameObject {
    constructor(playground) {
        super();
        this.playground = playground;
        this.ctx = this.playground.game_map.ctx;

        this.state = null; // win：胜利 lose：失败
        this.win_img = new Image();
        this.win_img.src = "https://cdn.acwing.com/media/article/image/2021/12/17/1_8f58341a5e-win.png";

        this.lose_img = new Image();
        this.lose_img.src = "https://cdn.acwing.com/media/article/image/2021/12/17/1_9254b5f95e-lose.png";
    }

    start() {

    }

    add_listening_events() {
        let outer = this;
        let $canvas = this.playground.game_map.$canvas;

        $canvas.on('click', function () {
            outer.playground.hide();
            outer.playground.root.menu.show();
        })
    }

    win() {
        let outer = this;
        this.state = "win";
        setTimeout(function () {
            outer.add_listening_events();
        }, 1000);
    }

    lose() {
        let outer = this;
        this.state = "lose";
        setTimeout(function () {
            outer.add_listening_events();
        }, 1000);
    }

    late_update() {
        this.render();
    }

    render() {
        let height = this.playground.height / 2;
        let width = this.playground.width / 2;
        if(this.state === "win") {
            this.ctx.drawImage(this.win_img, width - height / 2, height - height / 2, height, height);
        } else if(this.state === "lose") {
            this.ctx.drawImage(this.lose_img, width - height / 2, height - height / 2, height, height);
        }
    }
}