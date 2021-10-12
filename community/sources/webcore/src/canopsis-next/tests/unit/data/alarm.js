import Faker from 'faker';

import { fakeMeta } from '@unit/data/request-data';
import { fakeTimestamp } from './date';

export const fakeAlarm = () => ({
  _id: Faker.datatype.string(),
  t: fakeTimestamp(),
  entity: {
    _id: Faker.datatype.string(),
    name: Faker.datatype.string(),
    impact: [Faker.datatype.string()],
    depends: [Faker.datatype.string()],
    enable_history: [fakeTimestamp()],
    measurements: null,
    enabled: Faker.datatype.boolean(),
    infos: {
      criticity: {
        name: Faker.datatype.string(),
        description: Faker.datatype.string(),
        value: Faker.datatype.string(),
      },
    },
    type: 'resource',
    component: Faker.datatype.string(),
  },
  v: {
    state: {
      _t: 'stateinc',
      t: fakeTimestamp(),
      a: Faker.datatype.string(),
      m: Faker.datatype.string(),
      val: 3,
    },
    status: {
      _t: 'statusinc',
      t: fakeTimestamp(),
      a: Faker.datatype.string(),
      m: Faker.datatype.string(),
      val: 1,
    },
    component: Faker.datatype.string(),
    connector: Faker.datatype.string(),
    connector_name: Faker.datatype.string(),
    creation_date: fakeTimestamp(),
    activation_date: fakeTimestamp(),
    display_name: Faker.datatype.string(),
    initial_output: Faker.datatype.string(),
    output: Faker.datatype.string(),
    initial_long_output: Faker.datatype.string(),
    long_output: Faker.datatype.string(),
    long_output_history: [Faker.datatype.string()],
    last_update_date: fakeTimestamp(),
    last_event_date: fakeTimestamp(),
    resource: Faker.datatype.string(),
    tags: [],
    parents: [],
    children: [],
    total_state_changes: 1,
    extra: {},
    infos_rule_version: {},
    duration: 270,
    current_state_duration: 270,
    infos: {},
  },
  infos: {},
  links: {},
});

export const fakeAlarms = ({ count, limit = 10 } = {}) => Faker.datatype.array(count % limit).map(fakeAlarm);

export const fakeAlarmsResponse = ({ count, limit = 10, page = 1 } = {}) => ({
  data: fakeAlarms({ count, limit }),
  meta: fakeMeta({ count, limit, page }),
});
