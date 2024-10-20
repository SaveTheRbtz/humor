import { FaTwitter } from 'react-icons/fa';
import './Footer.css';

function Footer() {
  return (
    <footer className="footer">
      <a
        href="https://twitter.com/altsoph"
        target="_blank"
        rel="noopener noreferrer"
        className="footer-link"
      >
        <FaTwitter size={24} style={{ marginRight: '8px', marginLeft: '8px' }} />
        <b>@altsoph</b>
      </a>
      <a
        href="https://twitter.com/SaveTheRbtz"
        target="_blank"
        rel="noopener noreferrer"
        className="footer-link"
      >
        <FaTwitter size={24} style={{ marginRight: '8px', marginLeft: '8px' }} />
        <b>@SaveTheRbtz</b>
      </a>
    </footer>
  );
}

export default Footer;