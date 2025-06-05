import { StatusPopup } from './ui/components/StatusPopup';
import { renderReactComponent } from './utils';

// Event listener for when the DOM content is loaded
document.addEventListener('DOMContentLoaded', () => {
  const root = document.getElementById('root');
  if (root) {
    renderReactComponent(<StatusPopup />, root);
  }
});