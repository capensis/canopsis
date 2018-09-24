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
        view: {},
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
    createRow(row) {
      const { rows } = this.normalizedEntities.view[this.view._id];

      this.$set(this.normalizedEntities[rowSchema.key], row._id, row);
      this.$set(this.normalizedEntities[viewSchema.key][this.view._id], 'rows', [...rows, row._id]);
    },

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

        const oldRowId = this.config.rowId;
        const newRowId = this.settings.rowId;

        this.$set(this.normalizedEntities[widgetSchema.key], widget._id, widget);

        if (oldRowId !== newRowId) {
          if (oldRowId) {
            const oldRowWidgets = get(this.normalizedEntities, `${rowSchema.key}.${oldRowId}.widgets`, []);
            const filteredOldRowWidgets = oldRowWidgets.filter(oldWidget => oldWidget !== widget._id);

            this.$set(this.normalizedEntities[rowSchema.key][oldRowId], 'widgets', filteredOldRowWidgets);
          }

          const newRowWidgets = get(this.normalizedEntities, `${rowSchema.key}.${newRowId}.widgets`, []);
          const filteredNewRowWidgets = newRowWidgets.filter(v => v !== widget._id);

          this.$set(this.normalizedEntities[rowSchema.key][newRowId], 'widgets', [...filteredNewRowWidgets, widget._id]);
        }

        const view = denormalize(this.view._id, viewSchema, this.normalizedEntities);

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
