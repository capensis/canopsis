import { denormalize } from 'normalizr';

import queryMixin from '@/mixins/query';
import sideBarMixins from '@/mixins/side-bar/side-bar';
import entitiesViewMixin from '@/mixins/entities/view/view';
import entitiesUserPreferenceMixin from '@/mixins/entities/user-preference';

import { viewSchema } from '@/store/schemas';

import { convertUserPreferenceToQuery, convertWidgetToQuery } from '@/helpers/query';
import get from 'lodash/get';

export default {
  props: {
    config: {
      type: Object,
      required: true,
    },
  },
  mixins: [
    queryMixin,
    sideBarMixins,
    entitiesUserPreferenceMixin,
    entitiesViewMixin,
  ],
  data() {
    return {
      clonedView: {

      },
    };
  },
  computed: {
    widget() {
      return this.config.widget;
    },
  },
  methods: {
    isFormValid() {
      if (this.$validator) {
        return this.$validator.validateAll();
      }

      return true;
    },

    prepareSettingsWidget() {
      return this.settings.widget;
    },

    async submit() {
      const isFormValid = await this.isFormValid();

      if (isFormValid) {
        const widget = {
          ...this.widget,
          ...this.prepareSettingsWidget(),
        };

        const userPreference = {
          ...this.userPreference,
          widget_preferences: {
            ...this.userPreference.widget_preferences,
            ...this.settings.widget_preferences,
          },
        };

        const oldRowId = get(widget, '_embedded.parentId');
        const newRowId = this.rowId;

        this.entities.widget[widget._id] = widget;
        this.entities.viewRow[oldRowId].widgets = this.entities.viewRow[oldRowId].widgets.filter(v => v !== widget._id);
        this.entities.viewRow[newRowId].widgets =
          [...this.entities.viewRow[newRowId].widgets.filter(v => v !== widget._id), widget._id];

        const view = denormalize(this.view._id, viewSchema, this.entities);

        await Promise.all([
          this.createUserPreference({ userPreference }),
          this.updateView({ view }),
        ]);

        this.mergeQuery({
          id: widget._id,
          query: {
            ...convertWidgetToQuery(widget),
            ...convertUserPreferenceToQuery(userPreference),
          },
        });

        this.hideSideBar();
      }
    },
  },
};
