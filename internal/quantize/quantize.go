package quantize

import (
	"math"
)

// Quantize:

// quantized = (original - minValue) × scale

// Dequantize:

// original ≈ (quantized / scale) + minValue

type QuantizedVector struct {
    Values   []uint8
    MinValue float32
    Scale    float32
}



func Quantize(vector []float32) QuantizedVector {
    length := len(vector)

    if length == 0 {
        return QuantizedVector{}
    }

    minValue := vector[0]

    for _, value := range vector[1:] {
        if value < minValue {
            minValue = value
        }
    }

    maxValue := float32(1.0)
    scale := float32(255.0) / (maxValue - minValue)

    quantizedList := make([]uint8, length)

    for i := 0; i < length; i++ {
        val := vector[i]

        if val < minValue {
            val = minValue
        } else if val > maxValue {
            val = maxValue
        }

        scaled := (val - minValue) * scale

        quantizedList[i] = uint8(math.Round(float64(scaled)))
    }

    return QuantizedVector{
        Values:   quantizedList,
        MinValue: minValue,
        Scale:    scale,
    }
}

func Dequantize(vector QuantizedVector) []float32 {
    length := len(vector.Values)

    dequantizedList := make([]float32, length)

    for i := 0; i < length; i++ {
        val := vector.Values[i]

        descaled := (float32(val) / vector.Scale) + vector.MinValue

        dequantizedList[i] = descaled
    }

    return dequantizedList
}