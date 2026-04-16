import AxiosMockAdapter from 'axios-mock-adapter';
import Faker from 'faker';

import { API_ROUTES } from '@/config';

import request from '@/services/request';

import llmModule from '@/store/modules/entities/llm';

describe('Entities llm module', () => {
  const axiosMockAdapter = new AxiosMockAdapter(request);

  beforeEach(() => {
    axiosMockAdapter.reset();
  });

  afterEach(() => {
    jest.clearAllMocks();
  });

  test('Fetch LLM history without store. Action: fetchLlmHistoryWithoutStore', async () => {
    const id = Faker.datatype.uuid();
    const params = {
      user: Faker.internet.userName(),
      only_off_topic: true,
      page: 1,
      limit: 10,
    };
    const body = { data: [], meta: { total_count: 0 } };

    axiosMockAdapter
      .onGet(`${API_ROUTES.llms.list}/${id}/history`, { params })
      .reply(200, body);

    const result = await llmModule.actions.fetchLlmHistoryWithoutStore({}, { id, params });

    expect(result).toEqual(body);
  });

  test('Fetch LLM chats without store. Action: fetchLlmChatsWithoutStore', async () => {
    const id = Faker.datatype.uuid();
    const params = { user: Faker.internet.userName() };
    const body = { data: [] };

    axiosMockAdapter
      .onGet(`${API_ROUTES.llms.list}/${id}/chats`, { params })
      .reply(200, body);

    const result = await llmModule.actions.fetchLlmChatsWithoutStore({}, { id, params });

    expect(result).toEqual(body);
  });

  test('Fetch LLM users without store. Action: fetchLlmUsersWithoutStore', async () => {
    const id = Faker.datatype.uuid();
    const body = { data: [Faker.internet.userName()] };

    axiosMockAdapter
      .onGet(`${API_ROUTES.llms.list}/${id}/users`)
      .reply(200, body);

    const result = await llmModule.actions.fetchLlmUsersWithoutStore({}, { id });

    expect(result).toEqual(body);
  });

  test('Fetch LLM messages without store. Action: fetchLlmMessagesWithoutStore', async () => {
    const id = Faker.datatype.uuid();
    const params = { chat: Faker.datatype.uuid() };
    const body = { data: [] };

    axiosMockAdapter
      .onGet(`${API_ROUTES.llms.list}/${id}/messages`, { params })
      .reply(200, body);

    const result = await llmModule.actions.fetchLlmMessagesWithoutStore({}, { id, params });

    expect(result).toEqual(body);
  });

  test('Link LLM history in bulk without store. Action: bulkLinkLlmHistory', async () => {
    const data = { items: [{ _id: Faker.datatype.uuid() }] };
    const body = { updated: 1 };

    axiosMockAdapter
      .onPost(API_ROUTES.llms.bulkHistoryLink, data)
      .reply(200, body);

    const result = await llmModule.actions.bulkLinkLlmHistory({}, { data });

    expect(result).toEqual(body);
  });
});
