import { omit, isEmpty, isNumber } from 'lodash';

import { PATTERNS_FIELDS } from '@/constants';

const PATTERN_FIELDS_LIST = Object.values(PATTERNS_FIELDS);

/**
 * Converts one LLM chat history API message into the structure used by AI chat components:
 * pattern payloads are grouped under `patterns`, server-side version indices are shifted to zero-based,
 * and `emptyPatterns` is set for the first version when the previous row lacks pattern fields.
 *
 * @param {Object} message
 * @param {Object} [prevMessage]
 * @returns {Object}
 */
export const llmChatHistoryServerMessageToMessage = (message, prevMessage) => {
  const newMessage = omit(message, PATTERN_FIELDS_LIST);

  const patterns = PATTERN_FIELDS_LIST.reduce((acc, field) => {
    if (!message[field]) {
      return acc;
    }

    acc[field] = message[field];

    return acc;
  }, {});

  if (!isEmpty(patterns)) {
    newMessage.patterns = patterns;
  }

  if (prevMessage && message.version === 1) {
    const emptyPatterns = PATTERN_FIELDS_LIST.some(field => !prevMessage[field]);
    newMessage.emptyPatterns = emptyPatterns;
  }

  if (isNumber(message.version)) {
    newMessage.version = message.version - 1;
  }

  if (isNumber(message.from_version)) {
    newMessage.from_version = message.from_version - 1;
  }

  return newMessage;
};

/**
 * Maps a page of LLM chat history API messages to UI messages, passing each row the previous row
 * for first-version `emptyPatterns` detection.
 *
 * @param {Object[]} [serverMessages=[]]
 * @returns {Object[]}
 */
export const llmChatHistoryServerMessagesToMessages = (serverMessages = []) => (
  serverMessages.map((serverMessage, index) => (
    llmChatHistoryServerMessageToMessage(serverMessage, serverMessages?.[index - 1])
  ))
);
