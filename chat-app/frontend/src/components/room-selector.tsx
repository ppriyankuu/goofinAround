import React from 'react';

interface RoomSelectorProps {
  onSelect: (room: string) => void;
}

const RoomSelector: React.FC<RoomSelectorProps> = ({ onSelect }) => {
  const rooms = ['general', 'technical', 'fun'];

  return (
    <div className="bg-white p-8 rounded-lg shadow-md w-full max-w-sm">
      <h2 className="text-2xl font-bold mb-4">Select a Room</h2>
      <div className="space-y-2">
        {rooms.map((room) => (
          <button
            key={room}
            onClick={() => onSelect(room)}
            className="w-full bg-blue-500 hover:bg-blue-700 text-white font-bold py-2 px-4 rounded focus:outline-none focus:shadow-outline"
          >
            {room.charAt(0).toUpperCase() + room.slice(1)}
          </button>
        ))}
      </div>
    </div>
  );
};

export default RoomSelector;