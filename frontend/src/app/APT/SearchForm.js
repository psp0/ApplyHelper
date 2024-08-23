import React, { useState } from 'react';
import { fetchAPTData } from '../../../lib/api';
import './searchForm.style.css';

const SearchForm = ({ setAptData, setLoading }) => {
  const [searchParams, setSearchParams] = useState({
    suplyAreaCode: '',
    houseNm: '',
  });

  const handleInputChange = (e) => {
    const { name, value } = e.target;
    setSearchParams({ ...searchParams, [name]: value });
  };

  const handleSearch = async () => {
    setLoading(true);
    try {
      const data = await fetchAPTData(1, searchParams);
      setAptData(data);
      setLoading(false);
    } catch (error) {
      alert('Failed to fetch data');
      setLoading(false);
    }
  };

  return (
    <div className="search-form">
      <select name="suplyAreaCode" value={searchParams.suplyAreaCode} onChange={handleInputChange}>
        <option value="">수도권</option>
        <option value="">공급지역 전체</option>
        <option value="서울">서울</option>
        <option value="경기">경기</option>
        <option value="인천">인천</option>
        <option value="부산">부산</option>
        <option value="대구">대구</option>
        <option value="대전">대전</option>
        <option value="광주">광주</option>
        <option value="울산">울산</option>
        <option value="세종">세종</option>
        <option value="강원">강원</option>
        <option value="충북">충북</option>
        <option value="충남">충남</option>
        <option value="전북">전북</option>
        <option value="전남">전남</option>
        <option value="경북">경북</option>
        <option value="경남">경남</option>
        <option value="제주">제주</option> 
      </select>
      <input
        type="text"
        name="houseNm"
        placeholder="주택명(시행사명)"
        value={searchParams.houseNm}
        onChange={handleInputChange}
      />
      <button type="button" onClick={handleSearch}>
        조회
      </button>
    </div>
  );
};

export default SearchForm;
