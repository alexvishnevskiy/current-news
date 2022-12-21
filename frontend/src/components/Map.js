import React, { useState } from 'react';
import WorldMap from "react-world-map";

export default function Map() {
  const [selected, onSelect] = useState(null);

  return (
    <div className="App">
      <WorldMap selected={selected} onSelect={onSelect} />
    </div>
  );
}