class MultiPlayerSocket {
    constructor(playground) {
        this.playground = playground;
        this.ws = new WebSocket("ws://localhost:8089/wss/multiplayer");

        this.start();
    }

    start() {
        this.receive();
    }

    receive() {
        let outer = this;
        this.ws.onmessage = function (e) {
            let data = JSON.parse(e.data);
            let {response, event} = data;
            if(event === 'login') {
                outer.receive_login(response.codeMsg);
            } else {
                if(response.data.uuid === outer.uuid) return false;
                if(event === 'create_player') {
                    outer.receive_create_player(response.data.uuid, response.data.username, response.data.photo);
                } else if(event === 'move_to'){
                    outer.receive_move_to(response.data.uuid, response.data.tx, response.data.ty);
                } else if(event === 'shoot_fireball') {
                    outer.receive_shoot_fireball(response.data.uuid, response.data.tx, response.data.ty, response.data.ball_uuid)
                } else if(event === 'attack') {
                    outer.receive_attack(response.data.uuid, response.data.attacked_uuid, response.data.x, response.data.y, response.data.angle, response.data.damage, response.data.ball_uuid)
                } else if(event === 'blink') {
                    outer.receive_bink(response.data.uuid, response.data.tx, response.data.ty)
                } else if(event === 'message') {
                    outer.receive_message(response.data.username, response.data.text);
                }
            }
        }
    }

    send_login(username, photo) {
        this.username = username;
        this.photo = photo;
        let outer = this;
        this.ws.send(JSON.stringify({
            'seq': outer.uuid + '-login',
            'event': 'login',
            'data': {
                'userId': outer.uuid,
                'appId': 101,
                'username': outer.username,
                'photo': outer.photo,
            }
        }));
        // 心跳
        const heartbeat = () => {
            console.log("定时心跳:" + outer.username);
            this.ws.send(JSON.stringify({
                'seq': outer.uuid + "-heartbeat",
                'event': "heartbeat",
                "data":{}
            }))
        }
        // 定时心跳
        setInterval(heartbeat, 30 * 1000)
    }

    receive_login(codeMsg) {
        let username = this.username;
        let photo = this.photo;
        if(codeMsg === 'success') {
            this.send_create_player(username, photo)
        }
    }

    send_create_player(username, photo) {
        let outer = this;
        this.ws.send(JSON.stringify({
            'seq': outer.uuid + "-create_player",
            'event': "create_player",
            'data': {
                'appId': 101,
                'uuid': outer.uuid,
                'username': username,
                'photo': photo
            },
        }));
    }

    receive_create_player(uuid, username, photo) {
        let player = new Player(
            this.playground,
            this.playground.width / 2 / this.playground.scale,
            0.5,
            0.05,
            "white",
            0.15,
            "enemy",
            username,
            photo,
        );
        player.uuid = uuid;
        this.playground.players.push(player);
    }

    get_player(uuid) {
        let players = this.playground.players;
        for(let i = 0; i < players.length; i ++) {
            let player = players[i];
            if(player.uuid === uuid) {
                return player;
            }
        }
        return null;
    }

    send_move_to(tx, ty) {
        let outer = this;
        this.ws.send(JSON.stringify({
            'seq': outer.uuid + "-move_to",
            'event': "move_to",
            'data': {
                'appId': 101,
                'uuid': outer.uuid,
                'tx': tx,
                'ty': ty,
            },
        }));
    }

    receive_move_to(uuid, tx, ty) {
        let player = this.get_player(uuid);
        if(player) {
            player.move_to(tx, ty);
        }
    }

    send_shoot_fireball(tx, ty, ball_uuid) {
        let outer = this;
        this.ws.send(JSON.stringify({
            'seq': outer.uuid + "-shoot_fireball",
            'event': "shoot_fireball",
            'data': {
                'appId': 101,
                'uuid': outer.uuid,
                'tx': tx,
                'ty': ty,
                'ball_uuid': ball_uuid,
            },
        }));
    }

    receive_shoot_fireball(uuid, tx, ty, ball_uuid) {
        let player = this.get_player(uuid);
        if(player) {
            let fireball = player.shoot_fireball(tx, ty);
            fireball.uuid = ball_uuid;
        }
    }

    send_attack(attacked_uuid, x, y, angle, damage, ball_uuid) {
        let outer = this;
        console.log(angle, damage);
        this.ws.send(JSON.stringify({
            'seq': outer.uuid + "-attack",
            'event': "attack",
            'data': {
                'appId': 101,
                'uuid': outer.uuid,
                'attacked_uuid': attacked_uuid,
                'x': x,
                'y': y,
                'angle': angle,
                'damage': damage,
                'ball_uuid': ball_uuid,
            },
        }));
    }

    receive_attack(uuid, attacked_uuid, x, y, angle, damage, ball_uuid) {
        let attacker = this.get_player(uuid);
        let attacked = this.get_player(attacked_uuid);
        if(attacker && attacked) {
            attacked.receive_attack(x, y, angle, damage, ball_uuid, attacker);
        }
    }

    send_blink(tx, ty) {
        let outer = this;
        this.ws.send(JSON.stringify({
            'seq': outer.uuid + "-blink",
            'event': 'blink',
            'data': {
                'appId': 101,
                'uuid': outer.uuid,
                'tx': tx,
                'ty': ty,
            },
        }));
    }

    receive_bink(uuid, tx, ty) {
        let player = this.get_player(uuid);
        if(player) {
            player.blink(tx, ty);
        }
    }

    send_message(username, text) {
        let outer = this;
        this.ws.send(JSON.stringify({
            'seq': outer.uuid + "-message",
            'event': 'message',
            'data': {
                'appId': 101,
                'uuid': outer.uuid,
                'username': username,
                'text': text,
            }
        }));
    }

    // 死了也可以说话
    receive_message(username, text) {
        this.playground.chat_field.add_message(username, text);
    }
}
