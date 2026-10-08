import { Handlebars } from '@/helpers/handlebars';

import { handlebarsTagsHelperMixin } from '@/mixins/widget/handlebars/tags-helper';

import AlarmsListTable from '@/components/widgets/alarm/partials/alarms-list-table.vue';

describe('handlebarsTagsHelperMixin', () => {
  test('Alarm list registers tags helper on its own', () => {
    expect(AlarmsListTable.mixins).toContain(handlebarsTagsHelperMixin);
  });

  test('Tags helper is still available after service weather owner is destroyed', async () => {
    handlebarsTagsHelperMixin.beforeCreate();
    handlebarsTagsHelperMixin.beforeCreate();

    try {
      handlebarsTagsHelperMixin.destroyed();

      const template = Handlebars.compile('{{tags}}');
      const result = await template({ alarm: { v: { tags: ['tag'] } } });

      expect(result).toContain('<c-alarm-tags-chips');
    } finally {
      handlebarsTagsHelperMixin.destroyed();
    }
  });
});
