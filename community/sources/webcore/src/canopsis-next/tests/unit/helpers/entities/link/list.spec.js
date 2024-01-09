import { harmonizeCategoryLinks, harmonizeLinks, harmonizeAlarmsLinks } from '@/helpers/entities/link/list';

describe('entities links list helper', () => {
  const category = 'test-category';
  const category2 = 'test-category2';
  const link = {
    icon_name: 'icon',
    label: 'label',
    url: 'url',
    rule_id: 'rule_id',
    actions: 'copy',
    single: false,
  };

  const link2 = {
    icon_name: 'icon',
    label: 'label2',
    url: 'url',
    rule_id: 'rule_id2',
    actions: 'copy',
    single: true,
  };

  const link3 = {
    icon_name: 'icon',
    label: 'label3',
    url: 'url',
    rule_id: 'rule_id',
    actions: 'copy',
    single: false,
  };

  const oldLink = {
    label: 'label3',
    link: 'link',
  };

  const links = {
    [category]: [link2, link],
    [category2]: [oldLink],
  };

  const alarm = {
    links: { [category2]: [oldLink] },
  };

  const alarm2 = {
    links: { [category]: [link2, link] },
  };

  const alarm3 = {
    links: { [category2]: [link3] },
  };

  const alarms = [alarm, alarm2];
  const alarms2 = [alarm2, alarm3];

  it('Correct harmonize category links correct arguments', () => {
    expect(harmonizeCategoryLinks(links, category))
      .toEqual([link, link2]);
  });

  it('Correct harmonize category links without arguments', () => {
    expect(harmonizeCategoryLinks())
      .toEqual([]);
  });

  it('Correct harmonize links with correct arguments', () => {
    expect(harmonizeLinks(links))
      .toEqual([link, link2, oldLink]);
  });

  it('Correct harmonize alarms links with alarms without intersection in links', () => {
    expect(harmonizeAlarmsLinks(alarms))
      .toEqual([]);
  });

  it('Correct harmonize alarms links with alarms with intersection in links', () => {
    expect(harmonizeAlarmsLinks(alarms2))
      .toEqual([link, link3]);
  });

  it('Correct harmonize alarms links with one alarm with only old links', () => {
    expect(harmonizeAlarmsLinks([alarm]))
      .toEqual([]);
  });

  it('Correct harmonize alarms links with one alarm', () => {
    expect(harmonizeAlarmsLinks([alarm2]))
      .toEqual([link]);
  });

  it('Correct harmonize links without arguments', () => {
    expect(harmonizeLinks())
      .toEqual([]);
  });
});
