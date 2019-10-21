import { normalize, denormalize } from 'normalizr';

import queryMixin from '@/mixins/query';
import sideBarMixins from '@/mixins/side-bar/side-bar';
import entitiesViewMixin from '@/mixins/entities/view';
import entitiesUserPreferenceMixin from '@/mixins/entities/user-preference';

import { viewSchema, viewTabSchema, widgetSchema } from '@/store/schemas';

import { convertUserPreferenceToQuery, convertWidgetToQuery } from '@/helpers/query';

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
    entitiesViewMixin,
    entitiesUserPreferenceMixin,
  ],
  data() {
    return {
      tabId: this.config.tabId,
      normalizedEntities: {
        [viewSchema.key]: {},
        [viewTabSchema.key]: {},
        [widgetSchema.key]: {},
      },
    };
  },
  computed: {
    activeView() {
      return this.config.viewId ? this.getViewById(this.config.viewId) : this.view;
    },

    widget() {
      return this.config.widget;
    },

    localTab() {
      return denormalize(this.tabId, viewTabSchema, this.normalizedEntities);
    },

  },
  mounted() {
    const { entities } = normalize(this.activeView, viewSchema);

    this.normalizedEntities = entities;
  },
  methods: {
    updateNormalizedEntity(key, entity) {
      this.$set(
        this.normalizedEntities,
        key,
        { ...(this.normalizedEntities[key] || {}), [entity._id]: entity },
      );
    },

    isFormValid() {
      if (this.$validator) {
        return this.$validator.validateAll();
      }

      return true;
    },

    prepareWidgetSettings() {
      return this.settings.widget;
    },

    prepareWidgetQuery(newQuery) {
      return newQuery;
    },

    async submit() {
      const isFormValid = await this.isFormValid();

      if (isFormValid) {
        const widget = {
          ...this.settings.widget,
          ...this.prepareWidgetSettings(),
        };

        const userPreference = {
          ...this.userPreference,
          widget_preferences: {
            ...this.userPreference.widget_preferences,
            ...this.settings.widget_preferences,
          },
        };

        /**
         * Put widget into local normalized store
         */

        this.updateNormalizedEntity(widgetSchema.key, widget);

        const view = denormalize(this.activeView._id, viewSchema, this.normalizedEntities);

        await Promise.all([
          this.createUserPreference({ userPreference }),
          this.updateView({ id: this.activeView._id, data: view }),
        ]);

        const oldQuery = this.getQueryById(this.widget._id);
        const newQuery = {
          ...convertWidgetToQuery(widget),
          ...convertUserPreferenceToQuery(userPreference),
        };

        this.updateQuery({
          id: widget._id,
          query: this.prepareWidgetQuery(newQuery, oldQuery),
        });

        this.hideSideBar();
      }
    },
  },
};
