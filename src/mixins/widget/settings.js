import get from 'lodash/get';
import { normalize, denormalize } from 'normalizr';

import queryMixin from '@/mixins/query';
import sideBarMixins from '@/mixins/side-bar/side-bar';
import entitiesViewMixin from '@/mixins/entities/view';
import entitiesUserPreferenceMixin from '@/mixins/entities/user-preference';

import { WIDGET_MIN_SIZE, WIDGET_MAX_SIZE } from '@/constants';
import { viewSchema, rowSchema, widgetSchema } from '@/store/schemas';

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
      normalizedEntities: {
        [viewSchema.key]: {},
        [rowSchema.key]: {},
        [widgetSchema.key]: {},
      },
    };
  },
  computed: {
    widget() {
      return this.config.widget;
    },

    localView() {
      return denormalize(this.view._id, viewSchema, this.normalizedEntities);
    },

    availableRows() {
      if (!this.localView) {
        return [];
      }

      return this.localView.rows.map((row) => {
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
    const { entities } = normalize(this.view, viewSchema);

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

    createRow(row) {
      const view = this.normalizedEntities.view[this.view._id];

      this.updateNormalizedEntity(rowSchema.key, row);
      this.updateNormalizedEntity(viewSchema.key, {
        ...view,
        rows: [...view.rows, row._id],
      });
    },

    isFormValid() {
      if (this.$validator) {
        return this.$validator.validateAll();
      }

      return true;
    },

    prepareSettingsWidget() {
      return this.widget;
    },

    prepareWidgetQuery(newQuery) {
      return newQuery;
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

        const oldRowId = this.config.rowId;
        const newRowId = this.settings.rowId;

        /**
         * Put widget into local normalized store
         */

        this.updateNormalizedEntity(widgetSchema.key, widget);

        if (oldRowId !== newRowId) {
          if (oldRowId) {
            const oldRow = get(this.normalizedEntities, `${rowSchema.key}.${oldRowId}`, { widgets: [] });

            /**
             * Remove widget from old row in local normalized store
             */
            this.updateNormalizedEntity(rowSchema.key, {
              ...oldRow,
              widgets: oldRow.widgets.filter(oldWidget => oldWidget !== widget._id),
            });
          }

          const newRow = get(this.normalizedEntities, `${rowSchema.key}.${newRowId}`, { widgets: [] });

          /**
           * Put widget widget into new row in local normalized store
           */
          this.updateNormalizedEntity(rowSchema.key, {
            ...newRow,
            widgets: [
              ...newRow.widgets.filter(oldWidget => oldWidget !== widget._id),
              widget._id,
            ],
          });
        }

        const view = denormalize(this.view._id, viewSchema, this.normalizedEntities);

        await Promise.all([
          this.createUserPreference({ userPreference }),
          this.updateView({ id: this.view._id, data: view }),
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
