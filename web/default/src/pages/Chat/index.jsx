import React, { useState, useEffect, useRef } from 'react';
import {
  Card,
  Form,
  Button,
  List,
  TextArea,
  Input,
  Icon,
  Popup
} from 'semantic-ui-react';
import axios from 'axios';

const ChatApp = () => {
  // 当前对话记录，每条消息格式：{ role: 'user' | 'bot', content: string, isTyping?: boolean }
  const [messages, setMessages] = useState([]);
  // 输入框内容
  const [input, setInput] = useState('');
  // 请求状态
  const [loading, setLoading] = useState(false);
  // 侧边栏是否可见
  const [sidebarVisible, setSidebarVisible] = useState(false);
  // 历史聊天记录，数组中每项：{ id, name, messages, date }
  const [chatHistory, setChatHistory] = useState([]);
  // 搜索历史记录的关键词
  const [searchTerm, setSearchTerm] = useState('');
  // 输入框初始位置，desktop 第一次进入时居中，之后下移到底部
  const [inputPosition, setInputPosition] = useState('center');
  // 是否为移动端
  const [isMobile, setIsMobile] = useState(window.innerWidth < 768);

  // 聊天容器滚动到底部
  const chatEndRef = useRef(null);
  useEffect(() => {
    if (chatEndRef.current) {
      chatEndRef.current.scrollIntoView({ behavior: 'smooth' });
    }
  }, [messages]);

  // 监听屏幕宽度变化，实现响应式调整
  useEffect(() => {
    const handleResize = () => {
      setIsMobile(window.innerWidth < 768);
    };
    window.addEventListener('resize', handleResize);
    return () => window.removeEventListener('resize', handleResize);
  }, []);

  // 加载后端保存的历史聊天记录
  useEffect(() => {
    const loadChatHistory = async () => {
      try {
        const response = await axios.get('/api/chathistory');
        // 假设返回格式为 { history: [ { id, name, messages, date }, ... ] }
        setChatHistory(response.data.history);
      } catch (error) {
        console.error('加载历史记录失败：', error);
      }
    };
    loadChatHistory();
  }, []);

  // 新聊天：保存当前聊天到前后端历史记录，并清空当前对话
  const startNewChat = async () => {
    if (messages.length > 0) {
      const newHistoryItem = {
        id: Date.now(),
        name: messages[0].content.substring(0, 7),
        messages: messages,
        date: new Date()
      };
      // 更新前端状态
      setChatHistory([newHistoryItem, ...chatHistory]);
      // 调用后端接口保存聊天记录
      try {
        await axios.post('/api/saveChat', newHistoryItem);
      } catch (error) {
        console.error('保存聊天记录失败：', error);
      }
      setMessages([]);
    }
    // 重置输入框位置（桌面首次进入时在中间）
    setInputPosition('center');
  };

  // 发送消息及处理回复（包含逐字显示效果），并将消息保存到后端
  const sendMessage = async () => {
    const trimmedInput = input.trim();
    if (!trimmedInput) return;
    // 追加用户消息
    const userMessage = { role: 'user', content: trimmedInput };
    setMessages(prev => [...prev, userMessage]);
    setInput('');
    setLoading(true);
    // 如果是首次对话，将输入框位置下移到底部
    if (inputPosition === 'center') {
      setInputPosition('bottom');
    }
    // 保存用户消息到后端（可选：记录每条消息）
    try {
      await axios.post('/api/message', { message: trimmedInput, role: 'user' });
    } catch (error) {
      console.error('用户消息保存失败：', error);
    }
    try {
      // 调用后端接口获取回复（确保接口地址与返回数据格式正确）
      const response = await axios.post('/api/chat', { message: trimmedInput });
      const reply = response.data.reply;
      // 模拟逐字输出回复内容
      let currentIndex = 0;
      let currentReply = '';
      const interval = setInterval(async () => {
        currentReply += reply[currentIndex];
        setMessages(prev => {
          const newMessages = [...prev];
          const lastMsg = newMessages[newMessages.length - 1];
          if (lastMsg && lastMsg.role === 'bot' && lastMsg.isTyping) {
            lastMsg.content = currentReply;
          } else {
            newMessages.push({ role: 'bot', content: currentReply, isTyping: true });
          }
          return newMessages;
        });
        currentIndex++;
        if (currentIndex >= reply.length) {
          clearInterval(interval);
          // 移除 isTyping 标记
          setMessages(prev => {
            const newMessages = [...prev];
            if (newMessages[newMessages.length - 1].isTyping) {
              newMessages[newMessages.length - 1].isTyping = false;
            }
            return newMessages;
          });
          // 保存回复到后端
          try {
            await axios.post('/api/message', { message: reply, role: 'bot' });
          } catch (error) {
            console.error('回复消息保存失败：', error);
          }
        }
      }, 50); // 每50毫秒输出1个字符，可根据需要调整
    } catch (error) {
      console.error('请求失败：', error);
      const errorMessage = { role: 'bot', content: '请求错误，请稍后重试。' };
      setMessages(prev => [...prev, errorMessage]);
    }
    setLoading(false);
    // 发送后主动让输入框失去焦点，适用于移动端关闭软键盘
    document.activeElement.blur();
  };

  // 筛选符合搜索条件的历史记录
  const filteredHistory = chatHistory.filter(item =>
    item.name.toLowerCase().includes(searchTerm.toLowerCase())
  );

  // 删除历史记录项
  const deleteHistoryItem = async id => {
    setChatHistory(chatHistory.filter(item => item.id !== id));
    // 后端删除（根据接口要求调整）
    try {
      await axios.delete(`/api/chathistory/${id}`);
    } catch (error) {
      console.error('删除历史记录失败：', error);
    }
  };

  // 重命名历史记录项
  const renameHistoryItem = async id => {
    const newName = prompt('请输入新的聊天名称');
    if (newName) {
      setChatHistory(
        chatHistory.map(item =>
          item.id === id ? { ...item, name: newName } : item
        )
      );
      // 调用后端更新接口
      try {
        await axios.put(`/api/chathistory/${id}`, { name: newName });
      } catch (error) {
        console.error('重命名失败：', error);
      }
    }
  };

  // 根据屏幕尺寸调整侧边栏宽度
  const sidebarWidth = isMobile ? '80%' : '25%';

  return (
    <div
      style={{
        position: 'relative',
        width: '100vw',
        height: '100vh',
        overflow: 'hidden',
        background: '#f5f5f5'
      }}
    >
      {/* 侧边栏：覆盖效果 */}
      {sidebarVisible && (
        <div
          style={{
            position: 'absolute',
            top: 0,
            right: 0,
            width: sidebarWidth,
            height: '100%',
            background: '#fff',
            boxShadow: '-2px 0 5px rgba(0,0,0,0.1)',
            borderTopLeftRadius: '10px',
            borderBottomLeftRadius: '10px',
            display: 'flex',
            flexDirection: 'column',
            padding: '10px',
            zIndex: 1000
          }}
        >
          {/* 上半部分：关闭按钮、说明、搜索 */}
          <div style={{ flex: '0 0 auto', position: 'relative' }}>
            <Popup
              content="关闭侧边栏"
              trigger={
                <Button
                  circular
                  icon="close"
                  onClick={() => setSidebarVisible(false)}
                  style={{ position: 'absolute', top: '10px', right: '10px' }}
                />
              }
            />
            <div style={{ marginTop: '40px', marginBottom: '10px', textAlign: 'center' }}>
              <div>这里是说明文字</div>
            </div>
            <Input
              icon="search"
              placeholder="搜索聊天记录..."
              value={searchTerm}
              onChange={e => setSearchTerm(e.target.value)}
              fluid
            />
          </div>
          {/* 下半部分：动态历史聊天记录 */}
          <div style={{ flex: 1, overflowY: 'auto', marginTop: '10px' }}>
            <List divided relaxed>
              {filteredHistory.map(item => (
                <List.Item key={item.id}>
                  <List.Content floated="right">
                    <Popup
                      content="删除"
                      trigger={
                        <Button
                          icon="trash"
                          circular
                          onClick={() => deleteHistoryItem(item.id)}
                        />
                      }
                    />
                    <Popup
                      content="重命名"
                      trigger={
                        <Button
                          icon="edit"
                          circular
                          onClick={() => renameHistoryItem(item.id)}
                        />
                      }
                    />
                  </List.Content>
                  <List.Content>
                    <List.Header>{item.name}</List.Header>
                    <List.Description>
                      {new Date(item.date).toLocaleString()}
                    </List.Description>
                  </List.Content>
                </List.Item>
              ))}
            </List>
          </div>
        </div>
      )}

      {/* 内置容器，不允许整体滚动 */}
      <div
        style={{
          position: 'absolute',
          top: 0,
          left: 0,
          width: '100%',
          height: '100%',
          display: 'flex',
          flexDirection: 'column'
        }}
      >
        {/* 顶部工具栏：新聊天按钮和侧边栏开关 */}
        <div
          style={{
            position: 'relative',
            padding: '10px',
            display: 'flex',
            justifyContent: 'space-between',
            alignItems: 'center'
          }}
        >
          <Popup
            content="新聊天"
            trigger={
              <Button circular icon="plus" onClick={startNewChat} />
            }
          />
          {!sidebarVisible && (
            <Popup
              content="打开侧边栏"
              trigger={
                <Button
                  circular
                  icon="sidebar"
                  onClick={() => setSidebarVisible(true)}
                />
              }
            />
          )}
        </div>

        {/* 对话框区域：占约3/4屏幕，可滚动 */}
        <div
          style={{
            flex: 1,
            overflowY: 'auto',
            padding: '10px',
            marginTop: inputPosition === 'center' ? '25%' : '0',
            marginBottom: '10px'
          }}
        >
          <List relaxed>
            {messages.map((msg, idx) => (
              <List.Item
                key={idx}
                style={{
                  textAlign: msg.role === 'user' ? 'right' : 'left',
                  marginBottom: '10px'
                }}
              >
                <List.Content>
                  <List.Header>
                    {msg.role === 'user' ? '我' : 'ChatGPT'}
                  </List.Header>
                  <List.Description
                    style={{
                      whiteSpace: 'pre-wrap',
                      wordWrap: 'break-word'
                    }}
                  >
                    {msg.content}
                  </List.Description>
                </List.Content>
              </List.Item>
            ))}
            <div ref={chatEndRef}></div>
          </List>
        </div>

        {/* 输入区域：包含文本输入、发送与表情按钮 */}
        <div style={{ padding: '10px', position: 'relative' }}>
          <Form
            onSubmit={e => {
              e.preventDefault();
              sendMessage();
            }}
          >
            <Form.Field>
              <TextArea
                placeholder="请输入消息..."
                value={input}
                onChange={e => setInput(e.target.value)}
                onKeyPress={e => {
                  if (e.key === 'Enter' && !e.shiftKey) {
                    e.preventDefault();
                    sendMessage();
                  }
                }}
                style={{
                  borderRadius: '20px',
                  minHeight: '80px',
                  paddingRight: '80px'
                }}
              />
              {/* 右下角内置发送与表情按钮 */}
              <div
                style={{
                  position: 'absolute',
                  right: '20px',
                  bottom: '20px',
                  display: 'flex',
                  gap: '10px'
                }}
              >
                <Popup
                  content="表情"
                  trigger={
                    <Button
                      circular
                      icon="smile outline"
                      onClick={() => {
                        // 可添加表情选择器逻辑
                        alert('表情选择功能待实现');
                      }}
                    />
                  }
                />
                <Popup
                  content="发送"
                  trigger={
                    <Button
                      circular
                      icon="send"
                      primary
                      onClick={sendMessage}
                      loading={loading}
                    />
                  }
                />
              </div>
            </Form.Field>
          </Form>
        </div>
      </div>
    </div>
  );
};

export default ChatApp;
