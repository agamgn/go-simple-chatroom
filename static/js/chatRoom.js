$(function() {
    var $window = $(window);
    var $messageArea = $('#messageArea');       // 消息显示的区域
    var $inputArea = $('#inputArea');           // 输入消息的区域
    $inputArea.focus();
    var connected = false;                      // 用来判断是否连接的标志

    //随机颜色名字
    function getUsernameColor (username) {
        var COLORS = [
            '#e25412', '#54125f', '#d69851', '#a58421',
            '#854221', '#989658', '#cc4568', '#ddd521',
            '#eefe52', '#222dfd', '#abcd86', '#12ddfa'
        ];
        var hash = 8;
        for (var i = 0; i < username.length; i++) {
            hash = username.charCodeAt(i) + (hash << 5) - hash;
        }
        var index = Math.abs(hash % COLORS.length);
        return COLORS[index];
    }


    var nameColor = getUsernameColor( $("#name").text());
    $("#name").css('color', nameColor);


    //webSocket连接
    var socket = new WebSocket('ws://'+window.location.host+'/Room/WsRoom?name=' + $('#name').text());
    socket.onopen = function () {
        console.log("webSocket open");
        connected = true;
    };

    socket.onclose = function () {
        console.log("webSocket close");
        connected = false;
    };


    socket.onmessage = function(event) {
        var data = JSON.parse(event.data);
        console.log("revice:" , data);
        var name = data.name;
        var type = data.type;
        var msg = data.message;
        // type为0表示有人发消息
        var $messageDiv;
        if (type == 0) {
            var $usernameDiv = $('<span style="margin-right: 15px;font-weight: 700;overflow: hidden;text-align: right;"/>')
                    .text(name);
            if (name == $("#name").text()) {
                $usernameDiv.css('color', nameColor);
            } else {
                $usernameDiv.css('color', getUsernameColor(name));
            }
            var $messageBodyDiv = $('<span style="color: gray;"/>')
                    .text(msg);
            $messageDiv = $('<li style="list-style-type:none;font-size:25px;"/>')
                    .data('username', name)
                    .append($usernameDiv, $messageBodyDiv);
        }

        else {
            var $messageBodyDiv = $('<span style="color:#999999;">')
                    .text(msg);
            $messageDiv = $('<li style="list-style-type:none;font-size:15px;text-align:center;"/>')
                    .append($messageBodyDiv);
        }

        $messageArea.append($messageDiv);
        $messageArea[0].scrollTop = $messageArea[0].scrollHeight;   // 让屏幕滚动
    }

    $window.keydown(function (event) {
        if (event.which === 13) {
            sendMessage();
            typing = false;
        }
    });
    $("#sendBtn").click(function () {
        sendMessage();
    });


    function sendMessage () {
        var inputMessage = $inputArea.val();
        if (inputMessage && connected) {
            $inputArea.val('');
            socket.send(inputMessage);
            console.log("send message:" + inputMessage);
        }
    }
});
