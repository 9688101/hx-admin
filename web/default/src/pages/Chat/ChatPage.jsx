import React, { useState } from 'react';
import { Card, Form, Button, List, TextArea } from 'semantic-ui-react';
import axios from 'axios';

const ChatPage = () => {
  const [messages, setMessages] = useState([]); // 保存对话记录，格式：{ role: 'user'|'bot', content: string }
  const [input, setInput] = useState('');
  const [loading, setLoading] = useState(false);

  // 发送消息并请求后台接口
  const sendMessage = async () => {
    const trimmedInput = input.trim();
    if (!trimmedInput) return;

    // 将用户消息追加到对话记录中
    const userMessage = { role: 'user', content: trimmedInput };
    setMessages(prev => [...prev, userMessage]);
    setInput('');
    setLoading(true);

    try {
      // 发送请求到 Golang 后台
      const response = await axios.post('/api/chat', { message: trimmedInput });
      // 假设后台返回的数据格式为 { reply: '回复内容' }
      const botMessage = { role: 'bot', content: response.data.reply };
      setMessages(prev => [...prev, botMessage]);
    } catch (error) {
      console.error('请求失败：', error);
      const errorMessage = { role: 'bot', content: '请求错误，请稍后重试。' };
      setMessages(prev => [...prev, errorMessage]);
    }

    setLoading(false);
  };

  // 回车键触发发送消息（支持 Shift+Enter 换行）
  const handleKeyPress = (e) => {
    if (e.key === 'Enter' && !e.shiftKey) {
      e.preventDefault();
      sendMessage();
    }
  };

  return (
    <div className="chat-container" style={{ maxWidth: '800px', margin: '20px auto' }}>
      <Card fluid>
        <Card.Content header="ChatGPT Demo" />
        <Card.Content style={{ maxHeight: '500px', overflowY: 'auto' }}>
          <List relaxed divided>
            {messages.map((msg, idx) => (
              <List.Item key={idx} style={{ textAlign: msg.role === 'user' ? 'right' : 'left' }}>
                <List.Content>
                  <List.Header>{msg.role === 'user' ? '我' : 'ChatGPT'}</List.Header>
                  <List.Description style={{ whiteSpace: 'pre-wrap' }}>{msg.content}</List.Description>
                </List.Content>
              </List.Item>
            ))}
          </List>
        </Card.Content>
        <Card.Content extra>
          <Form>
            <Form.Field>
              <TextArea
                placeholder="请输入消息..."
                value={input}
                onChange={(e) => setInput(e.target.value)}
                onKeyPress={handleKeyPress}
                style={{ minHeight: '80px' }}
              />
            </Form.Field>
            <Button primary fluid onClick={sendMessage} loading={loading}>
              发送
            </Button>
          </Form>
        </Card.Content>
      </Card>
    </div>
  );
};

export default Chat;
