import { normalize, denormalize } from 'normalizr';

import queryMixin from '@/mixins/query';
import sideBarMixins from '@/mixins/side-bar/side-bar';
import entitiesViewMixin from '@/mixins/entities/view';
import entitiesUserPreferenceMixin from '@/mixins/entities/user-preference';

import { WIDGET_MIN_SIZE, WIDGET_MAX_SIZE } from '@/constants';
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

    availableRows() {
      if (!this.localTab) {
        return [];
      }

      return this.localTab.rows.map((row) => {
        const availableSize = row.widgets.reduce((acc, widget) => {
          if (widget._id !== this.widget._id) {
            acc.sm -= widget.size.sm;
            acc.md -= widget.size.md;
            acc.lg -= widget.size.lg;
          }

          return acc;
        }, { sm: WIDGET_MAX_SIZE, md: WIDGET_MAX_SIZE, lg: WIDGET_MAX_SIZE });

        return {
          _id: row._id,
          title: row.title,

          availableSize,
        };
      }).filter(({ availableSize }) =>
        availableSize.sm >= WIDGET_MIN_SIZE &&
        availableSize.md >= WIDGET_MIN_SIZE &&
        availableSize.lg >= WIDGET_MIN_SIZE);
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

        const view = denormalize(this.activeView._id, viewSchema, this.normalizedEntities);
        const tabIndex = view.tabs.findIndex(item => item._id === this.tabId);

        view.tabs[tabIndex].widgets.push(widget);
        // TODO: Default size value thanks to widget type + Compute default position
        view.tabs[tabIndex].layout.push({
          i: widget._id,
          x: 0,
          y: 0,
          w: 12,
          h: 3,
        });

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
