import React from 'react';
import Link from 'next/link';
import './menu.style.css';

const Menu = () => {
  return (
    <nav className="menu">
      <Link href="/" className="menu-button">Home</Link>
      <Link href="/APT" className="menu-button">APT</Link>
      <Link href="/RemndrAPT" className="menu-button">무순위/잔여세대</Link>
      <Link href="/contact" className="menu-button">오피스텔/생활숙박시설/도시형/(공공지원)민간임대</Link>
      <Link href="/RemndrAPT" className="menu-button">임의공급</Link>
    </nav>
  );
};

export default Menu;
