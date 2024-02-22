import Faker from 'faker';
import AxiosMockAdapter from 'axios-mock-adapter';

import { generateShallowRenderer } from '@unit/utils/vue';

import { API_ROUTES } from '@/config';

import request from '@/services/request';

import store from '@/store';

import { entitiesMapMixin } from '@/mixins/entities/map';

describe('Entities map mixin', () => {
  const axiosMockAdapter = new AxiosMockAdapter(request);

  const totalCount = 15;
  const maps = Array.from({ length: totalCount }).map(() => ({
    _id: Faker.datatype.string(),
    name: Faker.datatype.string(),
  }));
  const meta = {
    total_count: totalCount,
  };

  const factory = generateShallowRenderer({
    render() {},
    mixins: [entitiesMapMixin],
  }, { store });

  beforeEach(() => {
    axiosMockAdapter.reset();
  });

  test('Initial values applied', () => {
    const wrapper = factory();

    expect(wrapper.vm.mapsMeta).toEqual({});
    expect(wrapper.vm.mapsPending).toBe(false);
    expect(wrapper.vm.maps).toEqual([]);
  });

  test('Event maps list fetched', async () => {
    axiosMockAdapter
      .onGet(API_ROUTES.maps)
      .reply(200, { data: maps, meta });

    const wrapper = factory();

    const waitFetch = wrapper.vm.fetchMapsList();

    expect(wrapper.vm.mapsPending).toBe(true);

    await waitFetch;

    expect(wrapper.vm.mapsPending).toBe(false);

    expect(wrapper.vm.mapsMeta).toEqual(meta);
    expect(wrapper.vm.maps).toEqual(maps);
  });

  test('Event maps list fetched with params', async () => {
    const params = {
      page: Faker.datatype.number(),
    };
    axiosMockAdapter
      .onGet(API_ROUTES.maps, { params })
      .reply(200, { data: maps, meta });

    const wrapper = factory();

    await wrapper.vm.fetchMapsList({ params });

    expect(wrapper.vm.mapsMeta).toEqual(meta);
    expect(wrapper.vm.maps).toEqual(maps);
  });

  test('Event maps list fetched without store', async () => {
    const params = {
      page: Faker.datatype.number(),
    };
    axiosMockAdapter
      .onGet(API_ROUTES.maps, { params })
      .reply(200, { data: maps, meta });

    const wrapper = factory();

    const response = await wrapper.vm.fetchMapsListWithoutStore({ params });

    expect(response.data).toEqual(maps);
    expect(response.meta).toEqual(meta);
  });

  test('Maps list fetched with error', async () => {
    const error = Faker.datatype.string();
    const params = {
      page: Faker.datatype.number(),
    };
    axiosMockAdapter
      .onGet(API_ROUTES.maps, { params })
      .replyOnce(200, { data: maps, meta });

    const wrapper = factory();

    await wrapper.vm.fetchMapsList({ params });

    axiosMockAdapter
      .onGet(API_ROUTES.maps, { params })
      .reply(400, error);

    const originalError = console.error;
    console.error = jest.fn();

    try {
      await wrapper.vm.fetchMapsList({ params });
    } catch (err) {
      expect(err).toEqual(error);
    }

    expect(console.error).toBeCalledWith(error);

    expect(wrapper.vm.mapsPending).toBe(false);
    expect(wrapper.vm.mapsMeta).toEqual(meta);
    expect(wrapper.vm.maps).toEqual(maps);

    console.error = originalError;
  });

  test('Map fetched without store', async () => {
    const [map] = maps;
    const mapDetails = { ...map, parameters: {} };
    axiosMockAdapter
      .onGet(`${API_ROUTES.maps}/${map._id}`)
      .replyOnce(200, mapDetails);

    const wrapper = factory();

    const receivedMap = await wrapper.vm.fetchMapWithoutStore({ id: map._id });

    expect(receivedMap).toEqual(mapDetails);
  });

  test('Map state fetched without store', async () => {
    const [map] = maps;
    const mapDetails = { ...map, parameters: {} };
    axiosMockAdapter
      .onGet(`${API_ROUTES.mapState}/${map._id}`)
      .replyOnce(200, mapDetails);

    const wrapper = factory();

    const receivedMap = await wrapper.vm.fetchMapStateWithoutStore({ id: map._id });

    expect(receivedMap).toEqual(mapDetails);
  });

  test('Map created', async () => {
    const [map] = maps;
    axiosMockAdapter
      .onPost(API_ROUTES.maps)
      .replyOnce(200);

    const wrapper = factory();

    await wrapper.vm.createMap({ data: map });

    const [postRequestData] = axiosMockAdapter.history.post;

    expect(JSON.parse(postRequestData.data)).toEqual(map);
  });

  test('Map updated', async () => {
    const [map] = maps;
    axiosMockAdapter
      .onPut(`${API_ROUTES.maps}/${encodeURIComponent(map._id)}`)
      .replyOnce(200);

    const wrapper = factory();

    const updatedRule = {
      ...map,
      param: Faker.datatype.string(),
    };

    await wrapper.vm.updateMap({ data: updatedRule, id: map._id });

    const [putRequestData] = axiosMockAdapter.history.put;

    expect(JSON.parse(putRequestData.data)).toEqual(updatedRule);
  });

  test('Map removed', async () => {
    const [map] = maps;
    axiosMockAdapter
      .onDelete(`${API_ROUTES.maps}/${encodeURIComponent(map._id)}`)
      .replyOnce(200);

    const wrapper = factory();

    await wrapper.vm.removeMap({ id: map._id });

    expect(axiosMockAdapter.history.delete).toHaveLength(1);
  });

  test('Map removed', async () => {
    const [map] = maps;
    axiosMockAdapter
      .onDelete(API_ROUTES.bulkMaps)
      .replyOnce(200);

    const wrapper = factory();

    const data = [{ _id: map._id }];

    await wrapper.vm.bulkRemoveMaps({ data });

    expect(axiosMockAdapter.history.delete).toHaveLength(1);

    const [deleteRequest] = axiosMockAdapter.history.delete;
    expect(JSON.parse(deleteRequest.data)).toEqual(data);
  });
});
