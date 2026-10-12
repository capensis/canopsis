import { entitiesEntityInfoPropertyMixin } from '@/mixins/entities/entity-info-property';

/**
 * @TODO: remove this mixin in the future. Use entity-aliases-variables hook instead
 */
export const widgetColumnsEntityInfoPropertyMixin = {
  mixins: [entitiesEntityInfoPropertyMixin],
  mounted() {
    this.fetchEntityInfoPropertiesList({ params: { paginate: false } });
  },
  methods: {
    findAliasByColumnValue(columnValue, prefix) {
      const item = this.entityInfoProperties.find(property => (
        property.alias
        && columnValue === [prefix, 'infos', property.name, 'value'].filter(Boolean).join('.')
      ));

      return item?.alias;
    },
  },
};
