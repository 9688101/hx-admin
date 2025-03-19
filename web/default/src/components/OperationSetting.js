import React, { useEffect, useState } from 'react';
import { useTranslation } from 'react-i18next';
import { Form, Grid, Header, Button } from 'semantic-ui-react';
import { API, showError, showSuccess } from '../helpers';

const OperationSetting = () => {
  const { t } = useTranslation();
  const [inputs, setInputs] = useState({
    ChatLink: '',
  });
  const [originInputs, setOriginInputs] = useState({});
  const [loading, setLoading] = useState(false);

  const getOptions = async () => {
    const res = await API.get('/api/option/');
    const { success, message, data } = res.data;
    if (success) {
      const chatLinkOption = data.find(item => item.key === 'ChatLink') || { value: '' };
      setInputs({ ChatLink: chatLinkOption.value });
      setOriginInputs({ ChatLink: chatLinkOption.value });
    } else {
      showError(message);
    }
  };

  useEffect(() => {
    getOptions();
  }, []);

  const updateOption = async (key, value) => {
    setLoading(true);
    const res = await API.put('/api/option/', { key, value });
    const { success, message } = res.data;
    if (success) {
      setInputs(prev => ({ ...prev, [key]: value }));
      setOriginInputs(prev => ({ ...prev, [key]: value }));
    } else {
      showError(message);
    }
    setLoading(false);
  };

  const handleInputChange = (e, { name, value }) => {
    setInputs(prev => ({ ...prev, [name]: value }));
  };

  const submitConfig = async () => {
    if (originInputs.ChatLink !== inputs.ChatLink) {
      await updateOption('ChatLink', inputs.ChatLink);
    }
  };

  return (
    <Grid columns={1}>
      <Grid.Column>
        <Form loading={loading}>
          <Header as='h3'>{t('setting.operation.general.title')}</Header>
          <Form.Group widths={4}>
            <Form.Input
              label={t('setting.operation.general.chat_link')}
              name='ChatLink'
              onChange={handleInputChange}
              value={inputs.ChatLink}
              placeholder={t('setting.operation.general.chat_link_placeholder')}
              width={16}
            />
          </Form.Group>
          <Button primary onClick={submitConfig}>
            {t('setting.operation.general.buttons.save')}
          </Button>
        </Form>
      </Grid.Column>
    </Grid>
  );
};

export default OperationSetting;