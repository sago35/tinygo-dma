# ./examples/adctc

This is an example of performing an ADC conversion triggered by TC0.  
The results are captured using DMA.
The descriptor is a Linked Descriptor that refers to itself.  
The result of ADC conversion is output to the serial port.  

## environment

| pin | information |
| -- | -- |
| A0  | ADC input |
