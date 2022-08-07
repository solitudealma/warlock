let AC_GAME_OBJECTS = [];

class AcGameObject {
    constructor() {
        AC_GAME_OBJECTS.push(this);
        this.has_called_start = false; // 标记是否执行过start函数
        this.timedelta = 0; // 当前帧距离上一帧的时间间隔
        this.uuid = this.create_uuid();
    }

    create_uuid() {
        let res = "";
        for(let i = 0; i < 8; i ++) {
            let x = Math.random() * 10 | 0;
            res += x;
        }
        return res;
    }

    start() { // 只会在第一帧执行

    }

    update() { // 在每一帧执行

    }

    late_update() { // 在每一帧的最后执行一次

    }

    beforeDestroy() { // 被删除之前执行一些操作（更新分数等等）
    }

    destroy() { // 删掉该物体
        this.beforeDestroy();

        for(let i = 0; i < AC_GAME_OBJECTS.length; i ++) {
            if(AC_GAME_OBJECTS[i] === this) {
                AC_GAME_OBJECTS.splice(i, 1);
                break;
            }
        }
    }
}

let last_timestamp;
let AC_GAME_ANIMATION = function(timestamp) {
    for(let i = 0; i < AC_GAME_OBJECTS.length; i ++) {
        let obj = AC_GAME_OBJECTS[i];
        if(!obj.has_called_start) {
            obj.start();
            obj.has_called_start = true;
        }
        else {
            obj.timedelta = timestamp - last_timestamp;
            obj.update();
        }
    }

    for(let i = 0; i < AC_GAME_OBJECTS.length; i ++) {
        let obj = AC_GAME_OBJECTS[i];
        obj.late_update();
    }
    last_timestamp = timestamp;

    requestAnimationFrame(AC_GAME_ANIMATION);
}

requestAnimationFrame(AC_GAME_ANIMATION);
