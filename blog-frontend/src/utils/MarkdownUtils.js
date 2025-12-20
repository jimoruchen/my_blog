/**
 * Markdown工具类
 * 用于处理Markdown编辑器的渲染和解析
 */

import { marked } from 'marked';
import DOMPurify from 'dompurify';

// 配置marked选项
marked.setOptions({
  breaks: true,        // 允许回车换行
  gfm: true,           // 启用GitHub风格Markdown
  headerIds: true,     // 为标题添加id
  mangle: false,       // 不转义HTML
  smartLists: true,    // 使用更智能的列表行为
  smartypants: true,   // 使用更智能的标点符号
  xhtml: false         // 不使用自闭合标签
});

export default {
  /**
   * 将Markdown文本转换为HTML
   * @param {String} markdown Markdown文本
   * @returns {String} 安全的HTML
   */
  renderMarkdown(markdown) {
    if (!markdown) return '';
    // 使用marked将Markdown转换为HTML，然后使用DOMPurify清理HTML以防止XSS攻击
    return DOMPurify.sanitize(marked(markdown));
  },

  /**
   * 获取编辑器光标位置的行号和列号
   * @param {HTMLTextAreaElement} textarea 文本域元素
   * @returns {Object} 包含行号和列号的对象
   */
  getCursorPosition(textarea) {
    const position = textarea.selectionStart;
    const value = textarea.value.substring(0, position);
    const lines = value.split('\n');
    const lineCount = lines.length;
    const columnCount = lines[lineCount - 1].length;
    
    return {
      line: lineCount,
      column: columnCount,
      position: position
    };
  },

  /**
   * 在光标位置插入文本
   * @param {HTMLTextAreaElement} textarea 文本域元素
   * @param {String} text 要插入的文本
   */
  insertTextAtCursor(textarea, text) {
    const startPos = textarea.selectionStart;
    const endPos = textarea.selectionEnd;
    const scrollTop = textarea.scrollTop;
    
    textarea.value = textarea.value.substring(0, startPos) + 
                     text + 
                     textarea.value.substring(endPos, textarea.value.length);
    
    // 将光标位置设置到插入文本之后
    textarea.selectionStart = startPos + text.length;
    textarea.selectionEnd = startPos + text.length;
    textarea.scrollTop = scrollTop;
    
    // 触发input事件，确保Vue能够检测到变化
    const event = new Event('input', { bubbles: true });
    textarea.dispatchEvent(event);
  },

  /**
   * 获取随机颜色
   * @param {Number} seed 种子值
   * @returns {String} 颜色代码
   */
  getRandomColor(seed) {
    const colors = [
      '#FF6B6B', '#4ECDC4', '#45B7D1', '#FFA5A5', '#A5D8FF',
      '#FFD166', '#06D6A0', '#118AB2', '#EF476F', '#073B4C'
    ];
    return colors[seed % colors.length];
  }
};