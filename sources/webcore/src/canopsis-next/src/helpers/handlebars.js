import Handlebars from 'handlebars';

import dateFilter from '@/filters/date';

Handlebars.registerHelper('timestamp', (date) => {
  if (date) {
    return dateFilter(date, 'long');
  }

  return '';
});

Handlebars.registerHelper('entities', ({ hash }) => {
  const entityNameField = hash.name || 'entity.name';

  return new Handlebars.SafeString(`
    <div class="mt-2" v-for="watcherEntity in watcherEntitiess">
      <watcher-entity 
      :entity="watcherEntity"
      :template="config.entityTemplate"
      entityNameField="${entityNameField}"
      @addEvent="addEventToQueue"></watcher-entity>
    </div>
  `);
});

export function compile(template, context) {
  const handleBarFunction = Handlebars.compile(template);

  return handleBarFunction(context);
}

export default {
  compile,
};
