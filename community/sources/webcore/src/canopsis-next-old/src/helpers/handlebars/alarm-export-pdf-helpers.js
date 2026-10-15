import { isEmpty } from 'lodash';
import Handlebars from 'handlebars';

import { ALARM_EXPORT_PDF_FIELDS } from '@/constants';

import { harmonizeLinks } from '@/helpers/entities/link/list';

/**
 * Create paragraph HTML element
 *
 * @param {string} [content = '']
 * @returns {HTMLParagraphElement}
 */
export const createParagraph = (content = '') => {
  const p = document.createElement('p');

  p.style.padding = '12px';
  p.style.margin = '0px';

  if (content) {
    p.innerHTML = content?.outerHTML ?? content;
  }

  return p;
};

/**
 * Create table cell HTML element
 *
 * @param {boolean} [lastRow]
 * @param {boolean} [lastCell]
 * @returns {HTMLTableCellElement}
 */
export const createTableCell = (lastRow, lastCell) => {
  const td = document.createElement('td');

  td.style.padding = '0px';
  td.style.width = '50%';
  td.style.backgroundColor = '#fff';

  if (!lastRow) {
    td.style.borderBottom = 'solid #f5f5f5 2px';
  }

  if (!lastCell) {
    td.style.borderRight = 'solid #f5f5f5 2px';
  }

  return td;
};

/**
 * Create table row HTML element
 *
 * @param {string} [field = '']
 * @param {string} [value = '']
 * @param [last]
 * @returns {HTMLTableRowElement}
 */
export const createTableRow = (field = '', value = '', last) => {
  const tr = document.createElement('tr');
  const fieldTd = createTableCell(last);
  const valueTd = createTableCell(last, true);

  fieldTd.appendChild(createParagraph(field));
  valueTd.appendChild(createParagraph(value));

  tr.appendChild(fieldTd);
  tr.appendChild(valueTd);

  return tr;
};

/**
 * Create table HTML element
 *
 * @returns {HTMLTableElement}
 */
export function createTable() {
  const table = document.createElement('table');

  table.style.width = '100%';
  table.style.borderCollapse = 'collapse';
  table.style.tableLayout = 'fixed';
  table.style.marginBottom = '2px';
  table.style.marginRight = '2px';

  return table;
}

/**
 * Convert infos field to html string
 *
 * @returns {string}
 */
export function infosHelper() {
  const { [ALARM_EXPORT_PDF_FIELDS.infos]: infos } = this;

  if (isEmpty(infos)) {
    return '';
  }

  const table = createTable();
  const tbody = document.createElement('tbody');

  Object.values(infos).forEach((info, rootIndex, rootArray) => {
    Object.entries(info ?? {}).forEach(([key, value], index, array) => {
      const lastRow = index === array.length - 1 && rootIndex === rootArray.length - 1;

      tbody.appendChild(createTableRow(key, value, lastRow));
    });
  });

  table.appendChild(tbody);

  return new Handlebars.SafeString(table.outerHTML);
}

/**
 * Convert pbehavior_info field to html string
 *
 * @returns {string}
 */
export function pbehaviorInfoHelper() {
  const { [ALARM_EXPORT_PDF_FIELDS.pbehaviorInfo]: pbehaviorInfo } = this;

  if (isEmpty(pbehaviorInfo)) {
    return '';
  }

  const table = createTable();
  const tbody = document.createElement('tbody');

  tbody.appendChild(createTableRow('Enter time', pbehaviorInfo.timestamp));
  tbody.appendChild(createTableRow('Name', pbehaviorInfo.name));
  tbody.appendChild(createTableRow('Type', pbehaviorInfo.type_name));
  tbody.appendChild(createTableRow('Reason', pbehaviorInfo.reason_name, true));
  table.appendChild(tbody);

  return new Handlebars.SafeString(table.outerHTML);
}
/**
 * Convert ticket_info field to html string
 *
 * @returns {string}
 */
export function ticketHelper() {
  const { [ALARM_EXPORT_PDF_FIELDS.ticket]: ticket } = this;

  if (isEmpty(ticket)) {
    return '';
  }

  const table = createTable();
  const tbody = document.createElement('tbody');
  const ticketDataArray = Object.entries(ticket.ticket_data ?? {});

  tbody.appendChild(createTableRow('Ticket ID', ticket.ticket));
  tbody.appendChild(createTableRow('Ticket URL', ticket.ticket_url, !ticketDataArray.length));

  ticketDataArray.forEach(([field, value], index) => {
    tbody.appendChild(createTableRow(field, value, index === ticketDataArray.length - 1));
  });

  table.appendChild(tbody);

  return new Handlebars.SafeString(table.outerHTML);
}
/**
 * Convert comments field to html string
 *
 * @returns {string}
 */
export function commentsHelper() {
  const { [ALARM_EXPORT_PDF_FIELDS.comments]: comments = [] } = this;

  if (!comments?.length) {
    return '';
  }

  const table = createTable();
  const tbody = document.createElement('tbody');

  comments.forEach((comment, index) => {
    const tr = document.createElement('tr');
    const header = createParagraph();
    const author = document.createElement('strong');
    const timestamp = document.createElement('span');

    author.innerText = comment.a;
    author.style.marginRight = '10px';

    timestamp.innerText = comment.t;

    header.appendChild(author);
    header.appendChild(timestamp);

    const message = createParagraph(comment.m);
    const td = createTableCell(index === this.comments.length - 1, true);

    td.appendChild(header);
    td.appendChild(message);
    tr.appendChild(td);
    tbody.appendChild(tr);
  });

  table.appendChild(tbody);

  return new Handlebars.SafeString(table.outerHTML);
}

/**
 * Convert tags field to html string
 *
 * @returns {string}
 */
export function tagsHelper() {
  const { [ALARM_EXPORT_PDF_FIELDS.tags]: tags = [] } = this;

  if (!tags?.length) {
    return '';
  }

  const div = document.createElement('div');

  tags.forEach(tag => div.appendChild(createParagraph(tag)));

  return new Handlebars.SafeString(div.outerHTML);
}

/**
 * Convert links field to html string
 *
 * @returns {string}
 */
export function linksHelper() {
  const { [ALARM_EXPORT_PDF_FIELDS.links]: links = [] } = this;

  if (isEmpty(links)) {
    return '';
  }

  const div = document.createElement('div');

  harmonizeLinks(links).forEach(link => div.appendChild(createParagraph(link.url)));

  return new Handlebars.SafeString(div.outerHTML);
}

export const createInstanceWithHelpers = () => {
  const instance = Handlebars.create();

  instance.registerHelper(ALARM_EXPORT_PDF_FIELDS.infos, infosHelper);
  instance.registerHelper(ALARM_EXPORT_PDF_FIELDS.pbehaviorInfo, pbehaviorInfoHelper);
  instance.registerHelper(ALARM_EXPORT_PDF_FIELDS.ticket, ticketHelper);
  instance.registerHelper(ALARM_EXPORT_PDF_FIELDS.comments, commentsHelper);
  instance.registerHelper(ALARM_EXPORT_PDF_FIELDS.tags, tagsHelper);
  instance.registerHelper(ALARM_EXPORT_PDF_FIELDS.links, linksHelper);

  return instance;
};
