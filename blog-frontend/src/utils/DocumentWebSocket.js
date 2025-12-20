/**
 * 文档WebSocket客户端工具类
 * 用于处理文档编辑器的WebSocket连接
 */

export default class DocumentWebSocket {
  constructor(documentId, token, callbacks) {
    this.documentId = documentId;
    this.token = token;
    this.callbacks = callbacks || {};
    this.socket = null;
    this.reconnectAttempts = 0;
    this.maxReconnectAttempts = 5;
    this.reconnectTimeout = null;
  }

  /**
   * 连接WebSocket
   */
  connect() {
    if (!this.token) {
      this.triggerCallback('error', '未登录，无法建立实时连接');
      return;
    }

    // 创建WebSocket连接
    const protocol = window.location.protocol === 'https:' ? 'wss:' : 'ws:';
    const wsUrl = `${protocol}//${window.location.host}/ws/document/${this.documentId}?token=${this.token}`;
    
    this.socket = new WebSocket(wsUrl);

    this.socket.onopen = () => {
      console.log('WebSocket连接已建立');
      this.reconnectAttempts = 0;
      this.triggerCallback('open');
    };

    this.socket.onmessage = (event) => {
      try {
        const data = JSON.parse(event.data);
        this.triggerCallback('message', data);
      } catch (error) {
        console.error('处理WebSocket消息失败', error);
      }
    };

    this.socket.onclose = (event) => {
      console.log('WebSocket连接已关闭', event);
      this.triggerCallback('close', event);
      
      // 如果不是正常关闭，尝试重连
      if (!event.wasClean && this.reconnectAttempts < this.maxReconnectAttempts) {
        this.reconnectTimeout = setTimeout(() => {
          this.reconnectAttempts++;
          this.connect();
        }, 3000 * Math.pow(2, this.reconnectAttempts)); // 指数退避重连
      }
    };

    this.socket.onerror = (error) => {
      console.error('WebSocket错误', error);
      this.triggerCallback('error', '实时连接发生错误');
    };
  }

  /**
   * 发送消息
   * @param {Object} data 消息数据
   */
  send(data) {
    if (!this.socket || this.socket.readyState !== WebSocket.OPEN) {
      this.triggerCallback('error', '连接未建立，无法发送消息');
      return false;
    }

    try {
      this.socket.send(JSON.stringify(data));
      return true;
    } catch (error) {
      console.error('发送消息失败', error);
      this.triggerCallback('error', '发送消息失败');
      return false;
    }
  }

  /**
   * 发送内容更新
   * @param {String} content 文档内容
   */
  sendContentUpdate(content) {
    return this.send({
      type: 'update',
      content: content
    });
  }

  /**
   * 发送光标位置
   * @param {Number} position 光标位置
   * @param {String} username 用户名
   * @param {String} color 光标颜色
   */
  sendCursorPosition(position, username, color) {
    return this.send({
      type: 'cursor',
      position: position,
      username: username,
      color: color
    });
  }

  /**
   * 关闭连接
   */
  close() {
    if (this.reconnectTimeout) {
      clearTimeout(this.reconnectTimeout);
      this.reconnectTimeout = null;
    }
    
    if (this.socket) {
      this.socket.close();
      this.socket = null;
    }
  }

  /**
   * 触发回调
   * @param {String} type 回调类型
   * @param {*} data 回调数据
   */
  triggerCallback(type, data) {
    if (this.callbacks[type] && typeof this.callbacks[type] === 'function') {
      this.callbacks[type](data);
    }
  }
}