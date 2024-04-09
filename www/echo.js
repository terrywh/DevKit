export function setup(term, sock) {
    return new echo(term, sock);
}

const LOCAL_INPUT = 1;
const REMOTE_ECHO = 2;
const DEFAULT_TEXT_ENCODER = new TextEncoder();

class echo {
    constructor(term, sock) {
        this.term = term;
        this.sock = sock;
        this.q = [];
        this.o = null;
        this.px = -1;
        this.py = -1;
        this.term.onData(this.onPayload.bind(this, LOCAL_INPUT));
        this.sock.addEventListener("message", (e) => {
            this.onPayload(REMOTE_ECHO, new Uint8Array(e.data));
        });
        setInterval(this.onSend.bind(this), 160);
    }
    onPayload(type, payload) {
        if (type == REMOTE_ECHO) {
            if (this.px != -1 || this.py != -1) {
                this.positionReset();
            }
            this.term.write(payload);
        } else { // LOCAL_INPUT
            if (this.px == -1 || this.py == -1) {
                this.positionSave();
            }
            this.term.write(payload);
            this.q.push(payload);
        }
    }

    positionReset() {
        const dx = this.px - this.term.buffer.active.cursorX;
        const dy = this.py - this.term.buffer.active.cursorY;
        move(this.term, dx, dy);
        this.px = -1;
        this.py = -1;
    }

    positionSave() {
        this.px = this.term.buffer.active.cursorX;
        this.py = this.term.buffer.active.cursorY;
    }

    onSend(payload) {
        if (this.q.length > 0) {
            const data = DEFAULT_TEXT_ENCODER.encode(this.q.join(""));
            this.q.splice(0, this.q.length);
            this.sock.send(data);
        }
    }
}

const ESC = '\u001B[';
function move(term, x, y) {
    let move = '';
    if (x < 0) {
        move += ESC + (-x) + 'D';
    } else if (x > 0) {
        move += ESC + x + 'C';
    }
    if (y < 0) {
        move += ESC + (-y) + 'A';
    } else if (y > 0) {
        move += ESC + y + 'B';
    }
    term.write(move);
}

