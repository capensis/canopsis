import Faker from 'faker';
import { createPinia, setActivePinia } from 'pinia';

import { VUETIFY_ANIMATION_DELAY } from '@/config';

import { useSidebarStore } from '@/stores/modules/sidebar';

describe('sidebar store', () => {
  let store;

  beforeEach(() => {
    setActivePinia(createPinia());
    store = useSidebarStore();
  });

  afterEach(() => {
    vi.useRealTimers();
  });

  const showSidebar = () => {
    const id = Faker.datatype.uuid();
    const name = Faker.datatype.string();
    const config = Faker.helpers.createTransaction();

    store.show({ id, name, config });

    return { id, name, config };
  };

  it('shows a sidebar', () => {
    const { id, name, config } = showSidebar();

    expect(store.allIds).toEqual([id]);
    expect(store.byId[id]).toEqual({ id, name, config, hidden: false });
    expect(store.sidebars).toEqual([{ id, name, config, hidden: false }]);
  });

  it('uses an empty config by default', () => {
    const id = Faker.datatype.uuid();

    store.show({ id, name: Faker.datatype.string() });

    expect(store.byId[id].config).toEqual({});
  });

  it('keeps the currently opened sidebar', () => {
    const { id } = showSidebar();
    const nextId = Faker.datatype.uuid();

    const result = store.show({ id: nextId, name: Faker.datatype.string() });

    expect(result).toBe(id);
    expect(store.byId[nextId]).toBeUndefined();
  });

  it('marks a sidebar hidden and removes it after the animation', () => {
    vi.useFakeTimers();
    const { id } = showSidebar();

    store.hide({ id });

    expect(store.byId[id].hidden).toBe(true);

    vi.advanceTimersByTime(VUETIFY_ANIMATION_DELAY);

    expect(store.allIds).toEqual([]);
    expect(store.byId[id]).toBeUndefined();
  });

  it('ignores an unknown sidebar id', () => {
    store.hide({ id: 'missing' });

    expect(store.allIds).toEqual([]);
  });

  it('requires an id when hiding a sidebar', () => {
    expect(() => store.hide()).toThrow('Missed required parameter');
  });

  it('updates sidebar config', () => {
    const { id, config } = showSidebar();
    const patch = { title: Faker.datatype.string() };

    store.updateConfig({ id, config: patch });

    expect(store.byId[id].config).toEqual({ ...config, ...patch });
  });

  it('removes hidden sidebars before showing a new one', () => {
    vi.useFakeTimers();
    const { id } = showSidebar();
    const nextId = Faker.datatype.uuid();

    store.hide({ id });
    store.show({ id: nextId, name: Faker.datatype.string() });

    expect(store.allIds).toEqual([nextId]);
    expect(store.byId[id]).toBeUndefined();
  });
});
