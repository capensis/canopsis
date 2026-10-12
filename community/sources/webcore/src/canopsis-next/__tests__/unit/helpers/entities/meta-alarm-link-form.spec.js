import { formToMetaAlarmLinkRequest } from '@/helpers/entities/meta-alarm/link/form';

describe('meta alarm link form helpers', () => {
  it('uses the id of an existing meta alarm', () => {
    const metaAlarm = {
      _id: 'meta-alarm-id',
      name: 'Existing meta alarm',
    };

    expect(formToMetaAlarmLinkRequest({ metaAlarm, infos: [] })).toEqual({
      id: metaAlarm._id,
      infos: [],
    });
  });

  it('uses the typed name for a new meta alarm', () => {
    const metaAlarm = {
      _id: 'New meta alarm',
      noData: true,
    };

    expect(formToMetaAlarmLinkRequest({ metaAlarm, infos: [] })).toEqual({
      name: metaAlarm._id,
      infos: [],
    });
  });
});
